// Package service реализует бизнес-логику приложения.
package service

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
)

// rooms{
//     room {
//         id,
//         users: [
//             {
//                 id,
//                 name,
//                 symbol
//             },
//             {
//                 id,
//                 name,
//                 symbol
//             },
//         ],
//         positions: [
//             {
//                id,
//                symbol,
//                 user_id
//             },
//             ...
//         ],
//     },
//     ...
// }

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectedUser представляет подключённого пользователя в комнате игры.
type ConnectedUser struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	Connection  *websocket.Conn `json:"connection"`
	IsConnected bool            `json:"is_connected"`
}

// SymbolPosition описывает занятую позицию на игровом поле.
type SymbolPosition struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}

// RoomServer представляет комнату с пользователями и игровым состоянием.
type RoomServer struct {
	ID         uint64            `json:"id"`
	Users      []*ConnectedUser  `json:"users"`
	Positions  []*SymbolPosition `json:"symbol_positions"`
	BorderSize uint64            `json:"border_size"`
	GameStatus string            `json:"game_status"`
}

// WSServer управляет всеми комнатами и обработкой WebSocket-соединений.
type WSServer struct {
	Rooms        map[uint64]*RoomServer `json:"rooms"`
	ScoreService *ScoreService
	Mu           sync.Mutex
}

// GameRequest представляет входящее сообщение от клиента.
type GameRequest struct {
	Action     string      `json:"action,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Password   string      `json:"password,omitempty"`
	BorderSize uint64      `json:"size,omitempty"`
	Symbol     string      `json:"symbol,omitempty"`
}

// GameReponse представляет ответ сервера клиенту.
type GameReponse struct {
	Action      string      `json:"action"`
	Data        interface{} `json:"data,omitempty"`
	BoarderSize uint64      `json:"size,omitempty"`
	Symbol      string      `json:"symbol,omitempty"`
	UserID      *uuid.UUID  `json:"user_id,omitempty"`
}

// NewWsServer создаёт новый экземпляр WSServer.
func NewWsServer(scoreService *ScoreService) *WSServer {
	return &WSServer{
		Rooms:        make(map[uint64]*RoomServer),
		ScoreService: scoreService,
	}
}

// GameLoop обрабатывает основной цикл игры для пользователя.
func (ws *WSServer) GameLoop(
	currentUser *common.User,
	room *common.RoomSessionResponse,
	conn *websocket.Conn,
) bool {
	if ws.isRoomFull(currentUser.ID, room.ID, conn) {
		return true
	}
	ws.addUser(currentUser, room, conn)
	_, p, err := conn.ReadMessage()
	if err != nil {
		slog.Error("ReadMessage error:", slog.String("error", err.Error()))
		return true
	}
	if len(p) == 0 {
		slog.Warn("Received empty message")
		return true
	}
	var request GameRequest
	bytesReader := bytes.NewReader(p)
	if err := json.NewDecoder(bytesReader).Decode(&request); err != nil {
		slog.Error(
			"[wss]GameRequest",
			slog.String("error", err.Error()),
		)
	}
	slog.Info(
		"[wss]GameRequest",
		slog.Any("data", request),
	)
	return ws.proccessCommand(
		currentUser,
		room,
		request,
		p,
		conn,
	)
}

// RefreshConnection обновляет WebSocket-соединение для пользователя в комнате.
func (ws *WSServer) RefreshConnection(currentUser *common.User, room *common.RoomSessionResponse, conn *websocket.Conn) {
	if currentUser == nil {
		slog.Error("currentUser  is nil")
		return
	}
	if room == nil {
		slog.Error("room is nil")
		return
	}
	if ws.isUserInRoom(currentUser.ID, room.ID) {
		for _, user := range ws.Rooms[room.ID].Users {
			if user.ID == currentUser.ID && user.Connection != conn {
				user.Connection = conn
				user.IsConnected = true
				break
			}
		}
	}
}

// CloseConnection закрывает соединение пользователя и помечает его как отключённого.
func (ws *WSServer) CloseConnection(roomID uint64, conn *websocket.Conn) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	room, exists := ws.Rooms[roomID]
	if !exists || room == nil {
		slog.Warn(
			"attempted to close connection for non-existent room",
			slog.Uint64("room_id", roomID),
		)
		return
	}

	for _, user := range room.Users {
		if user.Connection == conn {
			conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection is close"),
			)
			conn.Close()
			user.Connection = nil
			user.IsConnected = false
			break
		}
	}
}

// proccessCommand обрабатывает действия, отправленные клиентом.
func (ws *WSServer) proccessCommand(
	currentUser *common.User,
	room *common.RoomSessionResponse,
	request GameRequest,
	message []byte,
	conn *websocket.Conn,
) bool {
	switch request.Action {
	case stepAction:
		ws.handleStep(room, &request)
	case resetGameAction:
		ws.handleResetGame(room)
	case resizeAction:
		ws.handleBorderResize(
			currentUser.ID,
			room,
			&request,
			message,
		)
	case selectSymbolAction:
		ws.handleSelectSymbol(
			currentUser.ID,
			room,
			&request,
		)
	case gameEndAction:
		ws.handleGameEnd(
			currentUser.ID,
			room,
			&request,
		)
	case exitRoomAction:
		return ws.handleExitRoom(
			currentUser,
			room,
			conn,
		)
	case closeRoomAction:
		ws.handleCloseRoom(ws.Rooms[room.ID])
	case newConnectionToRoomAction:
		ws.handleNewConnection(
			currentUser.ID,
			room,
			&request,
			conn,
		)
	}
	return false
}
