// Package middleware содержит промежуточные обработчики HTTP запросов
package middleware

import "net/http"

// CorsMiddleware создает middleware для обработки CORS (Cross-Origin Resource Sharing) заголовков
// Middleware выполняет следующие действия:
//   - Разрешает запросы с любых источников (Access-Control-Allow-Origin: *)
//   - Разрешает стандартные HTTP методы: GET, POST, PUT, DELETE, OPTIONS
//   - Разрешает заголовки Content-Type и Authorization
//   - Для OPTIONS запросов (preflight) сразу возвращает ответ без передачи дальше по цепочке
//
// Параметры:
//   - next http.Handler: следующий обработчик в цепочке middleware
//
// Возвращает:
//   - http.Handler: middleware функцию, которая добавляет CORS заголовки к ответу
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
