// Package http_handler предоставляет HTTP обработчики для API RemoteDesktop.
package http_handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
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

func (h *SessionHandler) Eidt(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	id := chi.URLParam(r, "id")
	if id == "" {
		resp.Message = "ID is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	guacToken := r.Header.Get("Guacamole-Token")
	if guacToken == "" {
		resp.Message = "Guacamole-Token is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	data, err := h.service.EditConnection(id, guacToken)
	if err != nil {
		resp.ResponseWrite(w, r, http.StatusNotFound)
		return
	}
	resp.Data = data
	resp.ResponseWrite(w, r, http.StatusOK)
}

func (h *SessionHandler) StoreConnection(w http.ResponseWriter, r *http.Request) {
	var resp helper.Response
	guacToken := r.Header.Get("Guacamole-Token")
	if guacToken == "" {
		resp.Message = "Guacamole-Token is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	var form common.GuacamoleConnectionRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		resp.Message = "Invalid JSON"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err := validate.Struct(form)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		humanReadableErrors, err := helper.LocalizedValidationMessages(r.Context(), errs)
		if err != nil {
			slog.Error(fmt.Sprintf("Error localizing validation messages: %s", err.Error()))
			resp.ResponseWrite(w, r, http.StatusInternalServerError)
			return
		}
		resp.Errors = humanReadableErrors
		resp.ResponseWrite(w, r, http.StatusUnprocessableEntity)
		return
	}
	if err := h.service.CreateConnection(&form, guacToken); err != nil {
		slog.Error(fmt.Sprintf("Error creating new connection:  %s", err.Error()))
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}
	resp.ResponseWrite(w, r, http.StatusOK)
}

func (h *SessionHandler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать обработку PUT запроса
}

func (h *SessionHandler) RemoveConnection(w http.ResponseWriter, r *http.Request) {
	var resp helper.Response
	id := chi.URLParam(r, "id")

	if id == "" {
		resp.Message = "Connection ID is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}

	guacToken := r.Header.Get("Guacamole-Token")
	if guacToken == "" {
		resp.Message = "Guacamole-Token is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}

	if err := h.service.DestroyConnection(id, guacToken); err != nil {
		slog.Error(fmt.Sprintf("Error removing connection:  %s", err.Error()))
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}

	resp.ResponseWrite(w, r, http.StatusOK)
}
