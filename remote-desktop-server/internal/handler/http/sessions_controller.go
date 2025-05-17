// Package http_handler предоставляет HTTP обработчики для API RemoteDesktop.
package http_handler

import (
	"net/http"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
)

// SessionHandler обрабатывает HTTP запросы для аутентификации пользователей.
type SessionHandler struct {
	service service.SessionService
}

// NewSessionHandler создает новый экземпляр SessionHandler.
//
// Параметры:
//   - service: сервис аутентификации, реализующий бизнес-логику
//
// Возвращает:
//   - *SessionHandler: указатель на созданный обработчик
func NewSessionHandler(service service.SessionService) *SessionHandler {
	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку GET запроса
}

func (h *SessionHandler) AddNewConnection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку POST запроса
}
