// Package service реализует бизнес-логику приложения.
package service

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
)

// opositeSymbol возвращает противоположный символ для игры (X/O)
//
// Параметры:
//   - mainSymbol: исходный символ ("X" или "O")
//
// Возвращает:
//   - string: противоположный символ или пустую строку если входной символ невалиден
func opositeSymbol(mainSymbol string) string {
	opositeSymbol := ""
	if mainSymbol == "" {
		opositeSymbol = ""
	}

	if mainSymbol == "X" {
		opositeSymbol = "O"
	}
	if mainSymbol == "O" {
		opositeSymbol = "X"
	}

	return opositeSymbol
}

// createRoom создает новую игровую комнату если она не существует
//
// Параметры:
//   - roomID: ID создаваемой комнаты
//
// Действия:
//   - Инициализирует комнату с дефолтными значениями если ее не существует
func (ws *WSServer) createRoom(roomID uint64) {
	if ws.Rooms[roomID] == nil {
		ws.Rooms[roomID] = &RoomServer{
			ID:         roomID,
			Users:      make([]*ConnectedUser, 0),
			Positions:  make([]*SymbolPosition, 0),
			BorderSize: DEFAULT_BORDER_SIZE,
		}
	}
}

// addUser добавляет пользователя в комнату
//
// Параметры:
//   - currentUser: данные пользователя
//   - room: целевая комната
//   - conn: WebSocket соединение
//
// Особенности:
//   - Использует мьютекс для потокобезопасности
//   - Не добавляет пользователя если он уже в комнате
func (ws *WSServer) addUser(currentUser *common.User, room *common.RoomSessionResponse, conn *websocket.Conn) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	ws.createRoom(room.ID)

	if !ws.isUserInRoom(currentUser.ID, room.ID) {
		ws.Rooms[room.ID].Users = append(
			ws.Rooms[room.ID].Users,
			&ConnectedUser{
				ID:          currentUser.ID,
				Name:        currentUser.Name,
				Symbol:      "",
				Connection:  conn,
				IsConnected: true,
			},
		)
	}
}

// isUserInRoom проверяет наличие пользователя в комнате
//
// Параметры:
//   - userID: ID пользователя
//   - roomID: ID комнаты
//
// Возвращает:
//   - bool: true если пользователь найден в комнате
func (ws *WSServer) isUserInRoom(userID uuid.UUID, roomID uint64) bool {
	room, exists := ws.Rooms[roomID]
	if !exists {
		return false
	}
	for _, user := range room.Users {
		if user.ID == userID {
			return true
		}
	}
	return false
}

// isRoomFull проверяет заполненность комнаты
//
// Параметры:
//   - userID: ID проверяемого пользователя
//   - roomID: ID комнаты
//   - conn: WebSocket соединение для отправки ошибки
//
// Возвращает:
//   - bool: true если комната заполнена (2 игрока) и пользователь не является участником
//
// Дополнительно:
//   - Отправляет сообщение об ошибке если комната заполнена
func (ws *WSServer) isRoomFull(userID uuid.UUID, roomID uint64, conn *websocket.Conn) bool {
	if ws.Rooms[roomID] != nil && len(ws.Rooms[roomID].Users) == 2 && !ws.isUserInRoom(userID, roomID) {
		err := conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "room is full"),
		)
		if err != nil {
			slog.Error(
				"Error sending room full message:",
				slog.String("error", err.Error()),
			)
		}
		return true
	}
	return false
}

// setSecondUserSymbol устанавливает символ второму игроку
//
// Параметры:
//   - roomId: ID комнаты
//
// Логика:
//  1. Если первый игрок выбрал символ, второму назначается противоположный
//  2. Отправляет уведомление второму игроку о назначенном символе
func (ws *WSServer) setSecondUserSymbol(roomId uint64) {
	currentRoom, exists := ws.Rooms[roomId]
	if !exists {
		return
	}
	if len(currentRoom.Users) == 2 {
		var prev *ConnectedUser
		for _, user := range currentRoom.Users {
			if user.Symbol != "" {
				if prev != nil {
					prev.Symbol = opositeSymbol(user.Symbol)
				}
				break
			}
			prev = user
		}
	}
	firstPlayerSymbol := ""
	secondarySymbol := ""
	for _, user := range currentRoom.Users {
		if user.Symbol != "" {
			firstPlayerSymbol = user.Symbol
		}
		if user.Symbol == "" {
			secondarySymbol = opositeSymbol(firstPlayerSymbol)
			if secondarySymbol != "" {
				user.Symbol = secondarySymbol
				resp := &GameReponse{
					Action: syncSymbolAction,
					Symbol: secondarySymbol,
				}
				raw, err := json.Marshal(resp)
				if err == nil && user.Connection != nil {
					user.Connection.WriteMessage(
						websocket.TextMessage,
						raw,
					)
				}
			}
		}
	}
}

// broadcastMessageToOther рассылает сообщение всем игрокам комнаты кроме указанного
//
// Параметры:
//   - currentUserID: ID игрока, которому НЕ нужно отправлять сообщение
//   - room: целевая комната
//   - raw: сырое сообщение для отправки
//
// Особенности:
//   - Использует мьютекс для потокобезопасности
//   - Логирует ошибки отправки
func (ws *WSServer) broadcastMessageToOther(
	currentUserID uuid.UUID,
	room *common.RoomSessionResponse,
	raw []byte,
) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	roomData, ok := ws.Rooms[room.ID]
	if !ok || roomData == nil {
		log.Printf("broadcastToSymbolPosition: room %d not found", room.ID)
		return
	}
	for _, currentUser := range roomData.Users {
		if currentUser.ID != currentUserID {
			if currentUser.Connection != nil {
				if err := currentUser.Connection.WriteMessage(
					websocket.TextMessage,
					raw,
				); err != nil {
					log.Println("WriteMessage error:", err)
				}
			}
		}
	}
}

// broadcastMessageToAll рассылает сообщение всем игрокам комнаты
//
// Параметры:
//   - room: целевая комната
//   - raw: сырое сообщение для отправки
//
// Особенности:
//   - Использует мьютекс для потокобезопасности
//   - Логирует ошибки отправки
func (ws *WSServer) broadcastMessageToAll(
	room *common.RoomSessionResponse,
	raw []byte,
) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	roomData, ok := ws.Rooms[room.ID]
	if !ok || roomData == nil {
		log.Printf("broadcastToSymbolPosition: room %d not found", room.ID)
		return
	}
	for _, currentUser := range roomData.Users {
		if currentUser.Connection != nil {
			if err := currentUser.Connection.WriteMessage(
				websocket.TextMessage,
				raw,
			); err != nil {
				log.Println("WriteMessage error:", err)
			}
		}
	}
}
