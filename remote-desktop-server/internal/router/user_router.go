// Package router предоставляет функциональность для настройки маршрутизации HTTP запросов.
package router

import (
	"github.com/go-chi/chi"
)

// usersRouterGroup регистрирует маршруты для работы с пользователями
//
// Параметры:
//   - users: chi.Router - роутер для регистрации маршрутов пользователей
//   - dependencies: содержит обработчики запросов (UserHandler)
//
// Регистрируемые маршруты:
//
//	GET /current - получение информации о текущем пользователе
func usersRouterGroup(users chi.Router) {
	users.Get("/current", dependencies.UserHandler.GetCurrentUser)
}
