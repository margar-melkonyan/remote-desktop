// Package http_handler предоставляет HTTP обработчики для API игры "Крестики-нолики".
package http_handler

import (
	"net/http"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/helper"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
)

// ScoreHandler обрабатывает HTTP запросы, связанные с игровыми результатами (счетами) пользователей
type ScoreHandler struct {
	service service.ScoreService
}

// NewScoreHandler создает новый экземпляр обработчика результатов
//
// Параметры:
//   - service: сервис для работы с результатами игр
//
// Возвращает:
//   - *ScoreHandler: указатель на созданный обработчик
func NewScoreHandler(service service.ScoreService) *ScoreHandler {
	return &ScoreHandler{
		service: service,
	}
}

// GetCurrentUserScores возвращает результаты текущего авторизованного пользователя
//
// Возможные коды ответа:
//   - 200: успешное получение результатов
//   - 401: пользователь не авторизован
//   - 500: внутренняя ошибка сервера
func (h *ScoreHandler) GetCurrentUserScores(w http.ResponseWriter, r *http.Request) {
	resp := helper.Response{}
	scores, err := h.service.GetCurrentUserScores(r.Context())
	if err != nil {
		resp.ResponseWrite(w, r, http.StatusInternalServerError)
		return
	}
	resp.Data = scores
	resp.ResponseWrite(w, r, http.StatusOK)
}
