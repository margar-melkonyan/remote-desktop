// Package http_handler предоставляет HTTP обработчики для API игры "Крестики-нолики".
package http_handler

import (
	"net/http"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/helper"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
)

// UserHandler обрабатывает HTTP запросы, связанные с пользователями
type UserHandler struct {
	service service.UserService
}

// NewUserHandler создает новый экземпляр обработчика пользователей
//
// Параметры:
//   - service: сервис для работы с пользователями
//
// Возвращает:
//   - *UserHandler: указатель на созданный обработчик
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetCurrentUser возвращает информацию о текущем авторизованном пользователе
//
// Возможные коды ответа:
//   - 200: успешное получение данных пользователя
//   - 401: пользователь не авторизован
//   - 409: конфликт при получении данных (например, пользователь не найден)
//   - 500: внутренняя ошибка сервера
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	user, err := h.service.GetCurrentUser(r.Context())
	if err != nil {
		resp.Message = err.Error()
		resp.ResponseWrite(w, r, http.StatusConflict)
		return
	}
	resp.Data = user
	resp.ResponseWrite(w, r, http.StatusOK)
}
