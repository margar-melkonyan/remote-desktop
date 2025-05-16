// Package router предоставляет функциональность для настройки маршрутизации HTTP запросов.
package router

import (
	"github.com/go-chi/chi"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common/dependency"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/handler/middleware"
)

var dependencies *dependency.AppDependencies

// NewRouter создает и настраивает маршрутизатор приложения
//
// Параметры:
//   - deps: зависимости приложения (*dependency.AppDependencies)
//
// Возвращает:
//   - *chi.Mux: настроенный маршрутизатор
//
// Особенности:
//   - Добавляет middleware для CORS и логирования
//   - Организует маршруты в иерархическую структуру
//   - Разделяет публичные и приватные маршруты
func NewRouter(deps *dependency.AppDependencies) *chi.Mux {
	dependencies = deps
	route := chi.NewMux()
	route.Use(
		middleware.CorsMiddleware,
		middleware.Logger,
	)

	route.Route("/auth", authRouterGroup)
	route.Route("/api", func(api chi.Router) {
		// Публичный маршрут получения списка комнат
		api.Get("/v1/rooms", dependencies.RoomHandler.GetRooms(dependencies.WSServer))
		// Приватные маршруты (требуют аутентификации)
		api.Route("/v1", func(v1 chi.Router) {
			v1.Use(middleware.AuthMiddleware(deps)) // Middleware аутентификации

			// Группы маршрутов:
			v1.Route("/rooms", roomsRouterGroup)   // Управление комнатами
			v1.Route("/users", usersRouterGroup)   // Работа с пользователями
			v1.Route("/scores", scoresRouterGroup) // Управление результатами игр
		})
	})

	return route
}
