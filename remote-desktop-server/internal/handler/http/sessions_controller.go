// Package handler предоставляет HTTP обработчики для работы с подключениями Guacamole.
// Реализует RESTful API для управления подключениями (создание, чтение, обновление, удаление).
package http_handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/helper"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
)

// SessionHandler обрабатывает HTTP запросы для работы с подключениями.
type SessionHandler struct {
	service *service.SessionService
}

// NewSessionHandler создает новый экземпляр SessionHandler.
//
// Параметры:
//   - service: сервис для работы с подключениями
//
// Возвращает:
//   - *SessionHandler: указатель на созданный обработчик
func NewSessionHandler(service *service.SessionService) *SessionHandler {
	return &SessionHandler{service: service}
}

// Get возвращает список подключений.
// Поддерживает фильтрацию по протоколу через query параметр.
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

// Edit возвращает информацию о конкретном подключении.
func (h *SessionHandler) Edit(w http.ResponseWriter, r *http.Request) {
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

// StoreConnection создает новое подключение.
func (h *SessionHandler) StoreConnection(w http.ResponseWriter, r *http.Request) {
	var resp helper.Response
	guacToken := r.Header.Get("Guacamole-Token")
	if guacToken == "" {
		resp.Message = "Guacamole-Token is required"
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	var form common.GuacamoleConnectionRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Ограничение тела запроса 10MB
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
		slog.Error(fmt.Sprintf("Error creating new connection: %s", err.Error()))
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}
	resp.ResponseWrite(w, r, http.StatusOK)
}

// UpdateConnection обновляет существующее подключение.
func (h *SessionHandler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
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
	var form common.GuacamoleConnectionRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Ограничение тела запроса 10MB
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

	if err := h.service.UpdateConnection(id, &form, guacToken); err != nil {
		slog.Error(fmt.Sprintf("Error updating connection: %s", err.Error()))
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}

	resp.ResponseWrite(w, r, http.StatusOK)
}

// RemoveConnection удаляет подключение.
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
		slog.Error(fmt.Sprintf("Error removing connection: %s", err.Error()))
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}

	resp.ResponseWrite(w, r, http.StatusOK)
}
