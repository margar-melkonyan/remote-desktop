// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"golang.org/x/crypto/bcrypt"
)

// handleStep обрабатывает ход игрока
//
// Параметры:
//   - room: текущая игровая комната
//   - request: запрос с данными хода
//
// Действия:
//  1. Парсит данные о позиции и символе
//  2. Обновляет состояние комнаты
//  3. Рассылает обновленные позиции всем игрокам
//  4. Устанавливает следующий ход для противоположного символа
func (ws *WSServer) handleStep(room *common.RoomSessionResponse, request *GameRequest) {
	rawSymbolPosition, ok := request.Data.(map[string]interface{})
	if ok {
		currentRoom := ws.Rooms[room.ID]
		symbolPosition := &SymbolPosition{
			ID:     rawSymbolPosition["id"].(string),
			Symbol: rawSymbolPosition["symbol"].(string),
		}
		currentRoom.GameStatus = inProcessStatus
		currentRoom.Positions = append(currentRoom.Positions, symbolPosition)
		ws.jsonToAll(room, &GameReponse{
			Action: getPositionsAction,
			Data: map[string]interface{}{
				"positions": currentRoom.Positions,
			},
			Symbol: opositeSymbol(symbolPosition.Symbol),
		})
	}
}

// handleResetGame сбрасывает состояние игры в комнате
//
// Параметры:
//   - room: текущая игровая комната
//
// Действия:
//  1. Очищает все сделанные ходы
//  2. Уведомляет всех игроков о сбросе
func (ws *WSServer) handleResetGame(room *common.RoomSessionResponse) {
	ws.Rooms[room.ID].Positions = make([]*SymbolPosition, 0)
	response := &GameReponse{
		Action: resetGameAction,
	}
	ws.jsonToAll(room, response)
}

// handleBorderResize обрабатывает изменение размера игрового поля
//
// Параметры:
//   - currentUserID: ID текущего пользователя
//   - room: игровая комната
//   - request: запрос с новым размером
//   - message: исходное сообщение
//
// Особенности:
//   - Доступно только создателю комнаты
//   - Рассылает изменение другим игрокам
func (ws *WSServer) handleBorderResize(
	currentUserID uuid.UUID,
	room *common.RoomSessionResponse,
	request *GameRequest,
	message []byte,
) {
	if currentUserID == room.CreatorID {
		ws.Rooms[room.ID].BorderSize = request.BorderSize
		ws.broadcastMessageToOther(currentUserID, room, message)
	}
}

// handleSelectSymbol обрабатывает выбор символа игроком
//
// Параметры:
//   - currentUserID: ID текущего пользователя
//   - room: игровая комната
//   - request: запрос с выбранным символом
//
// Действия:
//  1. Назначает символы игрокам (X/O)
//  2. Уведомляет другого игрока о выборе
func (ws *WSServer) handleSelectSymbol(
	currentUserID uuid.UUID,
	room *common.RoomSessionResponse,
	request *GameRequest,
) {
	currentRoom := ws.Rooms[room.ID]
	for id, user := range currentRoom.Users {
		if user.Symbol != "" {
			continue
		}
		if user.ID == currentUserID {
			currentRoom.Users[id].Symbol = request.Symbol
		} else {
			currentRoom.Users[id].Symbol = opositeSymbol(request.Symbol)
		}
	}
	response := &GameReponse{
		Action: selectedSymbolAction,
		Symbol: opositeSymbol(request.Symbol),
	}
	ws.jsonToOther(currentUserID, room, response)
}

// handleNewConnection обрабатывает новое подключение к комнате
//
// Параметры:
//   - currentUserID: ID подключающегося пользователя
//   - room: игровая комната
//   - request: запрос подключения
//   - conn: WebSocket соединение
//
// Действия:
//  1. Проверяет пароль для приватных комнат
//  2. Инициализирует состояние комнаты
//  3. Рассылает текущее состояние новому игроку
//  4. Устанавливает символы игрокам
func (ws *WSServer) handleNewConnection(
	currentUserID uuid.UUID,
	room *common.RoomSessionResponse,
	request *GameRequest,
	conn *websocket.Conn,
) {
	currentRoom := ws.Rooms[room.ID]
	if room.IsPrivate != nil && room.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(room.Password), []byte(request.Password))
		if err != nil {
			conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "password is not valid"),
			)
			ws.CloseConnection(room.ID, conn)
		}
	}
	ws.jsonToAll(room, &GameReponse{
		Action: newConnectionToRoomAction,
		UserID: &currentUserID,
	})
	ws.jsonToAll(room, &GameReponse{
		Action:      resizeAction,
		BoarderSize: currentRoom.BorderSize,
	})
	currentRoom.GameStatus = chooseSymbolStatus
	ws.jsonToAll(room, &GameReponse{
		Action: chooseSymbolAction,
		UserID: &currentRoom.Users[0].ID,
	})
	currentPlayerStep := currentRoom.Users[0].Symbol
	if len(currentRoom.Positions) != 0 {
		currentPlayerStep = opositeSymbol(currentRoom.Positions[len(currentRoom.Positions)-1].Symbol)
	}
	ws.jsonToAll(room, &GameReponse{
		Action: getPositionsAction,
		Data: map[string]interface{}{
			"positions": currentRoom.Positions,
		},
		Symbol: currentPlayerStep,
	})
	ws.setSecondUserSymbol(room.ID)
}

// handleExitRoom обрабатывает выход игрока из комнаты
//
// Параметры:
//   - currentUser: выходящий пользователь
//   - room: игровая комната
//   - conn: WebSocket соединение
//
// Возвращает:
//   - bool: true если обработка завершена
//
// Действия:
//  1. Фиксирует результат игры (если игра шла)
//  2. Уведомляет оставшегося игрока
//  3. Сбрасывает состояние комнаты
//  4. Закрывает соединение
func (ws *WSServer) handleExitRoom(
	currentUser *common.User,
	room *common.RoomSessionResponse,
	conn *websocket.Conn,
) bool {
	slog.Info(
		"User  exiting room:",
		slog.String("user_id", currentUser.ID.String()),
	)

	currentRoom := ws.Rooms[room.ID]
	var versusPlayer *ConnectedUser
	for _, user := range currentRoom.Users {
		if currentUser.ID != user.ID {
			versusPlayer = user
			break
		}
	}

	if versusPlayer != nil {
		if currentRoom.GameStatus == inProcessStatus {
			ws.ScoreService.scoreRepo.Create(context.Background(), &common.Score{
				IsWon:    1,
				UserID:   versusPlayer.ID.String(),
				Nickname: currentUser.Name,
			})
			ws.ScoreService.scoreRepo.Create(context.Background(), &common.Score{
				IsWon:    0,
				UserID:   currentUser.ID.String(),
				Nickname: versusPlayer.Name,
			})
		}

		ws.jsonToOther(currentUser.ID, room, &GameReponse{
			Action: chooseSymbolAction,
			UserID: &versusPlayer.ID,
		})
		currentRoom.Positions = make([]*SymbolPosition, 0)
		ws.jsonToOther(currentUser.ID, room, &GameReponse{
			Action: getPositionsAction,
			Data: map[string]interface{}{
				"positions": currentRoom.Positions,
			},
		})
		currentRoom.GameStatus = chooseSymbolStatus
		ws.Mu.Lock()
		defer ws.Mu.Unlock()
		versusPlayer.Symbol = ""
		currentRoom.Users = []*ConnectedUser{
			versusPlayer,
		}
	}
	conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection is close"),
	)
	conn.Close()
	return true
}

// handleGameEnd обрабатывает завершение игры
//
// Параметры:
//   - currentUserID: ID текущего пользователя
//   - room: игровая комната
//   - request: запрос с результатом игры
//
// Действия:
//  1. Сохраняет результат игры
//  2. Устанавливает статус "игра завершена"
//  3. Уведомляет о завершении игры
func (ws *WSServer) handleGameEnd(
	currentUserID uuid.UUID,
	room *common.RoomSessionResponse,
	request *GameRequest,
) {
	slog.Info(
		"Game ended for user:",
		slog.String("user_id", currentUserID.String()),
	)
	rawGameEndJson, ok := request.Data.(map[string]interface{})
	if !ok {
		return
	}
	isWon, ok := rawGameEndJson["is_won"].(float64)
	if !ok {
		return
	}
	userID, ok := rawGameEndJson["user_id"].(string)
	if !ok {
		return
	}
	nickname, ok := rawGameEndJson["versus_player_nickname"].(string)
	if !ok {
		return
	}
	score := &common.Score{
		IsWon:    isWon,
		UserID:   userID,
		Nickname: nickname,
	}
	ws.Rooms[room.ID].GameStatus = gameEndStatus
	ws.ScoreService.scoreRepo.Create(context.Background(), score)
	ws.jsonToOther(
		currentUserID,
		room, &GameReponse{
			Action: restartGameAction,
		})
}

// handleCloseRoom полностью закрывает комнату
//
// Параметры:
//   - room: комната для закрытия
//
// Действия:
//  1. Закрывает все соединения в комнате
//  2. Удаляет комнату из списка активных
func (ws *WSServer) handleCloseRoom(
	room *RoomServer,
) {
	for _, user := range room.Users {
		ws.CloseConnection(room.ID, user.Connection)
	}

	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	delete(ws.Rooms, room.ID)
}

// jsonToAll рассылает JSON сообщение всем игрокам комнаты
func (ws *WSServer) jsonToAll(room *common.RoomSessionResponse, response *GameReponse) {
	raw, err := json.Marshal(response)
	if err == nil {
		ws.broadcastMessageToAll(room, raw)
	}
}

// jsonToOther рассылает JSON сообщение другим игрокам комнаты
func (ws *WSServer) jsonToOther(currentUserID uuid.UUID, room *common.RoomSessionResponse, response interface{}) {
	raw, err := json.Marshal(response)
	if err == nil {
		ws.broadcastMessageToOther(currentUserID, room, raw)
	}
}
