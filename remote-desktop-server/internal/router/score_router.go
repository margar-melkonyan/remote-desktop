// Package router предоставляет функциональность для настройки маршрутизации HTTP запросов.
package router

import (
	"github.com/go-chi/chi"
)

// scoresRouterGroup регистрирует маршруты для работы с игровыми результатами
//
// Параметры:
//   - scores: chi.Router - роутер для регистрации маршрутов результатов
//   - dependencies: содержит обработчики запросов (ScoreHandler)
//
// Регистрируемые маршруты:
//
//	GET / - получение результатов текущего пользователя
func scoresRouterGroup(scores chi.Router) {
	scores.Get("/", dependencies.ScoreHandler.GetCurrentUserScores)
}
