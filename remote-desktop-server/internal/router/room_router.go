// Package router предоставляет функциональность для настройки маршрутизации HTTP запросов.
package router

import (
	"github.com/go-chi/chi"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/handler/ws"
)

// roomsRouterGroup регистрирует маршруты для работы с игровыми комнатами
//
// Параметры:
//   - rooms: chi.Router - роутер для регистрации маршрутов комнат
//   - dependencies: содержит обработчики запросов и WebSocket сервер
//
// Регистрируемые маршруты:
//
//	POST / - создание новой игровой комнаты
//	GET /{id}/info - получение информации о комнате
//	GET /{id} - WebSocket подключение к комнате
//	GET /my - список комнат текущего пользователя
//	DELETE /{id} - удаление комнаты
func roomsRouterGroup(rooms chi.Router) {
	rooms.Post("/", dependencies.RoomHandler.CreateRoom)
	rooms.Get("/{id}/info", dependencies.RoomHandler.GetRoom(dependencies.WSServer))
	rooms.Get("/{id}", ws.EnterRoom(dependencies))
	rooms.Get("/my", dependencies.RoomHandler.GetMyRooms(dependencies.WSServer))
	rooms.Delete("/{id}", dependencies.RoomHandler.DestroyRoom)
}
