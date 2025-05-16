// Package middleware содержит промежуточные обработчики HTTP запросов
package middleware

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// statusRecorder - обертка для http.ResponseWriter, которая запоминает статус код ответа
// и поддерживает интерфейс http.Hijacker для работы с WebSocket и другими соединениями
type statusRecorder struct {
	statusCode          int // сохраненный HTTP статус код ответа
	http.ResponseWriter     // встроенный интерфейс ResponseWriter
}

// Hijack реализует метод интерфейса http.Hijacker для поддержки upgrade соединений
// (например, WebSocket). Если исходный ResponseWriter не поддерживает Hijacker,
// возвращает ошибку.
func (sr *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := sr.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support hijacking")
}

// WriteHeader перехватывает и сохраняет статус код перед вызовом оригинального метода
func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

// Logger создает middleware для логирования информации о HTTP запросах.
// Логирует:
//   - HTTP метод (GET, POST и т.д.)
//   - Статус код ответа
//   - URI запроса
//   - Время выполнения запроса
//
// Параметры:
//   - next http.Handler: следующий обработчик в цепочке middleware
//
// Возвращает:
//   - http.Handler: middleware функцию, которая логирует информацию о запросе
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{
			statusCode:     200,
			ResponseWriter: w,
		}
		start := time.Now()
		next.ServeHTTP(sr, r)
		slog.Info(
			"Request info",
			slog.String("method", r.Method),
			slog.Int("status", sr.statusCode),
			slog.String("uri", r.URL.String()),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
