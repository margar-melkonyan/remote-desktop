// Package http_handler предоставляет HTTP обработчики для API RemoteDesktop.
package http_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/margar-melkonyan/remote-desktop.git/internal/common"
	"github.com/margar-melkonyan/remote-desktop.git/internal/helper"
	"github.com/margar-melkonyan/remote-desktop.git/internal/service"
)

// AuthHandler обрабатывает HTTP запросы для аутентификации пользователей.
type AuthHandler struct {
	service service.AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler.
//
// Параметры:
//   - service: сервис аутентификации, реализующий бизнес-логику
//
// Возвращает:
//   - *AuthHandler: указатель на созданный обработчик
func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// SingIn обрабатывает запрос на вход пользователя.
//
// Параметры:
//   - w: http.ResponseWriter для записи ответа
//   - r: *http.Request с данными запроса
//
// Логика работы:
//  1. Проверяет Content-Type запроса
//  2. Читает и парсит JSON тело запроса (макс. 10MB)
//  3. Валидирует входные данные
//  4. При ошибках валидации возвращает локализованные сообщения
//  5. Вызывает сервис аутентификации
//  6. Возвращает JWT токен при успешной аутентификации
//
// Возможные коды ответа:
//   - 200: успешный вход, возвращает токен
//   - 400: ошибка парсинга JSON
//   - 422: ошибки валидации
//   - 409: конфликт (неверные учетные данные)
//   - 500: внутренняя ошибка сервера
func (h *AuthHandler) SingIn(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	if resp.IsValidMediaType(w, r) {
		return
	}
	var form common.AuthSignInRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		slog.Error("Error decoding JSON: " + err.Error())
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
	token, err := h.service.SignIn(r.Context(), form)
	if err != nil {
		resp.Message = err.Error()
		resp.ResponseWrite(w, r, http.StatusConflict)
		return
	}
	resp.Data = token
	resp.ResponseWrite(w, r, http.StatusOK)
}

// SignUp обрабатывает запрос на регистрацию пользователя.
//
// Параметры:
//   - w: http.ResponseWriter для записи ответа
//   - r: *http.Request с данными запроса
//
// Логика работы:
//  1. Проверяет Content-Type запроса
//  2. Читает и парсит JSON тело запроса (макс. 10MB)
//  3. Валидирует входные данные
//  4. При ошибках валидации возвращает локализованные сообщения
//  5. Вызывает сервис регистрации
//  6. Возвращает статус создания
//
// Возможные коды ответа:
//   - 200: успешная регистрация
//   - 400: ошибка парсинга JSON
//   - 422: ошибки валидации
//   - 409: конфликт (пользователь уже существует)
//   - 500: внутренняя ошибка сервера
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	if resp.IsValidMediaType(w, r) {
		return
	}
	var form common.AuthSignUpRequest
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		slog.Error("Error decoding JSON: " + err.Error())
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
	if err := h.service.SignUp(r.Context(), form); err != nil {
		resp.Message = err.Error()
		resp.ResponseWrite(w, r, http.StatusConflict)
		return
	}
	resp.Message = "Created!"
	resp.ResponseWrite(w, r, http.StatusOK)
}
