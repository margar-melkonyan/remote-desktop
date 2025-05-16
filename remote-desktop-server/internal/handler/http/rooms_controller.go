// Package http_handler предоставляет HTTP обработчики для API игры "Крестики-нолики".
package http_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/helper"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
)

// RoomHandler обрабатывает HTTP запросы для работы с игровыми комнатами.
type RoomHandler struct {
	service service.RoomService
}

// NewRoomHandler создает новый экземпляр RoomHandler.
//
// Параметры:
//   - service: сервис комнат, реализующий бизнес-логику
//
// Возвращает:
//   - *RoomHandler: указатель на созданный обработчик
func NewRoomHandler(service service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

// GetRooms возвращает обработчик для получения списка всех комнат.
//
// Параметры:
//   - ws: WebSocket сервер для получения актуального состояния комнат
//
// Возвращает:
//   - http.HandlerFunc: обработчик, который:
//     1. Получает список всех комнат через RoomService
//     2. Возвращает список в формате JSON
//     3. Всегда возвращает статус 200 OK
func (h *RoomHandler) GetRooms(ws *service.WSServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := helper.Response{}
		data := h.service.GetAll(r.Context(), ws)
		resp.Data = data
		resp.ResponseWrite(w, r, http.StatusOK)
	}
}

// GetMyRooms возвращает обработчик для получения списка комнат текущего пользователя.
//
// Параметры:
//   - ws: WebSocket сервер для получения актуального состояния комнат
//
// Возвращает:
//   - http.HandlerFunc: обработчик, который:
//     1. Получает список комнат текущего пользователя через RoomService
//     2. Возвращает список в формате JSON
//     3. Всегда возвращает статус 200 OK
func (h *RoomHandler) GetMyRooms(ws *service.WSServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := helper.Response{}
		data := h.service.GetAllMy(r.Context(), ws)
		resp.Data = data
		resp.ResponseWrite(w, r, http.StatusOK)
	}
}

// GetRoomInfo возвращает информацию о конкретной комнате.
// Вспомогательный метод, используемый другими обработчиками.
//
// Параметры:
//   - r: HTTP запрос
//   - ws: WebSocket сервер
//
// Возвращает:
//   - helper.Response: ответ с данными комнаты или ошибкой
func (h *RoomHandler) GetRoomInfo(r *http.Request, ws *service.WSServer) helper.Response {
	resp := helper.Response{}
	param := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		resp.Errors = err
		return resp
	}
	data, err := h.service.GetById(r.Context(), id, ws)
	if err != nil {
		resp.Errors = err.Error()
		return resp
	}
	resp.Data = data
	return resp
}

// GetRoom возвращает обработчик для получения информации о конкретной комнате.
//
// Параметры:
//   - ws: WebSocket сервер
//
// Возвращает:
//   - http.HandlerFunc: обработчик, который:
//     1. Извлекает ID комнаты из URL параметров
//     2. Получает информацию о комнате через RoomService
//     3. Возвращает данные комнаты или ошибку
//     4. Возвращает статус 200 OK при успехе
func (h *RoomHandler) GetRoom(ws *service.WSServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h.GetRoomInfo(r, ws)
		resp.ResponseWrite(w, r, http.StatusOK)
	}
}

// CreateRoom обрабатывает запрос на создание новой комнаты.
//
// Логика работы:
//  1. Проверяет Content-Type запроса
//  2. Читает и парсит JSON тело запроса (макс. 10MB)
//  3. Валидирует входные данные
//  4. При ошибках валидации возвращает локализованные сообщения
//  5. Создает комнату через RoomService
//
// Возможные коды ответа:
//   - 200: комната успешно создана
//   - 400: ошибка парсинга JSON
//   - 422: ошибки валидации
//   - 500: внутренняя ошибка сервера
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	if resp.IsValidMediaType(w, r) {
		return
	}
	var form common.RoomRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		slog.Error("Error decoding JSON: ", slog.String("error", err.Error()))
		resp.ResponseWrite(w, r, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(&form)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		humanReadableErrors, err := helper.LocalizedValidationMessages(
			r.Context(),
			errs,
		)
		if err != nil {
			slog.Error("Error localizing validation messages: " + err.Error())
			resp.ResponseWrite(w, r, http.StatusInternalServerError)
			return
		}
		resp.Errors = humanReadableErrors
		resp.ResponseWrite(w, r, http.StatusUnprocessableEntity)
		return
	}
	if err := h.service.Create(r.Context(), form); err != nil {
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}
	resp.Message = "Created!"
	resp.ResponseWrite(w, r, http.StatusOK)
}

// DestroyRoom обрабатывает запрос на удаление комнаты.
//
// Логика работы:
//  1. Извлекает ID комнаты из URL параметров
//  2. Пытается удалить комнату через RoomService
//
// Возможные коды ответа:
//   - 200: комната успешно удалена
//   - 404: неверный ID комнаты
//   - 204: комната не найдена (уже удалена)
func (h *RoomHandler) DestroyRoom(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	param := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		resp.ResponseWrite(w, r, http.StatusNotFound)
		return
	}
	if err := h.service.DeleteById(r.Context(), id); err != nil {
		resp.ResponseWrite(w, r, http.StatusNoContent)
		return
	}
	resp.ResponseWrite(w, r, http.StatusOK)
}
