// Package http_handler предоставляет HTTP обработчики для API RemoteDesktop.
package http_handler

import (
	"net/http"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/helper"
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
	resp := helper.Response{}
	var protocol string
	hasProtocol := r.URL.Query().Has("protocol")
	if hasProtocol {
		protocol = r.URL.Query().Get("protocol")
	} else {
		protocol = "all"
	}
	guacToken := r.Header.Get("Guacamole-Token")
	if guacToken == "" {
		resp.Message = "Guacamole-Token is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	data, err := h.service.GetSession(protocol, guacToken)
	if err != nil {
		resp.ResponseWrite(w, r, http.StatusNotFound)
		return
	}
	resp.Data = data
	resp.ResponseWrite(w, r, http.StatusOK)
}

func (h *SessionHandler) StoreConnection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку POST запроса
}

func (h *SessionHandler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку PUT запроса
}

func (h *SessionHandler) RemoveConnection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку DELETE запроса
}
