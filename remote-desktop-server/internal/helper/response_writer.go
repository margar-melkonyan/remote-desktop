// Package helper предоставляет вспомогательные функции и структуры для работы с HTTP-ответами.
// Основная функциональность:
//   - Стандартизированная форма ответов API
//   - Проверка заголовков запросов
//   - Сериализация JSON-ответов
package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// Response представляет стандартную структуру ответа API.
// Используется для унификации всех ответов сервера.
// Поля:
//   - Data: основные данные ответа (может быть опущено)
//   - Errors: ошибки валидации или бизнес-логики (может быть опущено)
//   - Message: текстовое сообщение (может быть опущено)
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ResponseWriter интерфейс для записи HTTP-ответов.
// Определяет стандартный метод ResponseWrite для отправки ответов.
type ResponseWriter interface {
	ResponseWrite(w http.ResponseWriter, r *http.Request, status int)
}

// ResponseWrite записывает стандартизированный JSON-ответ.
//
// Параметры:
//   - w: HTTP ResponseWriter для записи ответа
//   - r: HTTP запрос
//   - status: HTTP статус код ответа
//
// Действия:
//
//  1. Устанавливает Content-Type: application/json
//  2. Записывает HTTP статус
//  3. Сериализует ответ в JSON
//  4. В случае ошибки сериализации записывает сообщение об ошибке
func (response *Response) ResponseWrite(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	raw, err := json.Marshal(response)
	if err != nil {
		slog.Error(err.Error())
		w.Write([]byte("something went wrong"))
	}
	w.Write(raw)
}

// IsValidMediaType проверяет корректность Content-Type заголовка запроса.
//
// Параметры:
//   - w: HTTP ResponseWriter
//   - r: HTTP запрос
//
// Возвращает:
//   - true: если Content-Type отсутствует или не равен "application/json"
//   - false: если Content-Type корректен
//
// Действия:
//
//	При некорректном Content-Type:
//	1. Устанавливает сообщение об ошибке
//	2. Отправляет ответ со статусом 415 Unsupported Media Type
func (response *Response) IsValidMediaType(w http.ResponseWriter, r *http.Request) bool {
	contentType := strings.TrimSpace(r.Header.Get("Content-Type"))
	if contentType == "" || contentType != "application/json" {
		response.Message = "Not valid content-type"
		response.ResponseWrite(w, r, http.StatusUnsupportedMediaType)
		return true
	}
	return false
}
