// Package middleware содержит промежуточные обработчики HTTP запросов
package middleware

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/margar-melkonyan/remote-desktop.git/internal/common"
	"github.com/margar-melkonyan/remote-desktop.git/internal/common/dependency"
	"github.com/margar-melkonyan/remote-desktop.git/internal/helper"
	"github.com/margar-melkonyan/remote-desktop.git/internal/service"
)

// responseWriterWrapper оборачивает http.ResponseWriter для добавления
// поддержки интерфейса http.Hijacker
type responseWriterWrapper struct {
	http.ResponseWriter
}

// Hijack реализует метод интерфейса http.Hijacker
func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support hijacking")
}

// AuthMiddleware создает middleware для аутентификации запросов
// Параметры:
//   - dependency *dependency.AppDependencies: зависимости приложения, включая репозитории
//
// Возвращает:
//
//	func(next http.Handler) http.Handler: middleware функцию, которая:
//	   1. Проверяет наличие токена в заголовке Authorization или query параметре token
//	   2. Валидирует токен с помощью service.CheckTokenIsNotExpired
//	   3. Ищет пользователя в репозитории по email из токена
//	   4. При успешной аутентификации добавляет email и данные пользователя в контекст
//	   5. При ошибках возвращает HTTP 401 с соответствующим сообщением
func AuthMiddleware(dependency *dependency.AppDependencies) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" && r.URL.Query().Get("token") != "" {
				token = r.URL.Query().Get("token")
			}
			if token == "" {
				resp := helper.Response{}
				resp.Message = "You should be authorized!"
				resp.ResponseWrite(w, r, http.StatusUnauthorized)
				return
			}
			claims, err := service.CheckTokenIsNotExpired(token)
			if err != nil {
				resp := helper.Response{}
				resp.Message = err.Error()
				resp.ResponseWrite(w, r, http.StatusUnauthorized)
				return
			}
			user, err := dependency.GlobalRepositories.UserRepository.FindByEmail(r.Context(), claims.Sub.Email)
			if err != nil {
				resp := helper.Response{}
				resp.Message = err.Error()
				resp.ResponseWrite(w, r, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), common.USER_MAIL, claims.Sub.Email)
			ctx = context.WithValue(ctx, common.USER, user)
			r = r.WithContext(ctx)

			wrappedWriter := &responseWriterWrapper{w}
			next.ServeHTTP(wrappedWriter, r)
		})
	}
}
