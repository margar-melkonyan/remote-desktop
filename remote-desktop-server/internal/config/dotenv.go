package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

// mustLoadEnv загружает переменные окружения из файла .env.
// В случае ошибки завершает работу приложения с panic.
//
// Использует:
//   - godotenv для загрузки переменных окружения
//
// Особенности:
//   - Вызывается при инициализации приложения
//   - Критическая для работы приложения функция
//   - При отсутствии/недоступности .env файла вызывает panic
//
// Пример использования:
//
//	mustLoadEnv() // Загружает .env или завершает приложение
func mustLoadEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("can't load .env")
		panic(err)
	}
}
