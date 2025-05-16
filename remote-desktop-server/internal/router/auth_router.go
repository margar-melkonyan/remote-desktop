// Package router предоставляет функциональность для настройки маршрутизации HTTP запросов.
package router

import (
	"github.com/go-chi/chi"
)

// authRouterGroup регистрирует маршруты для аутентификации в указанном роутере.
//
// Параметры:
//   - auth: chi.Router - роутер для регистрации маршрутов аутентификации
//   - dependencies: содержит обработчики запросов (AuthHandler)
//
// Регистрируемые маршруты:
//
//	POST /sign-in - обработка входа пользователя
//	POST /sign-up - обработка регистрации нового пользователя
func authRouterGroup(auth chi.Router) {
	auth.Post("/sign-in", dependencies.AuthHandler.SingIn)
	auth.Post("/sign-up", dependencies.AuthHandler.SignUp)
}
