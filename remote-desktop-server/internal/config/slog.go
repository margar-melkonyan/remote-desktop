package config

import (
	"log/slog"
	"os"
)

// setNewDefaultLogger инициализирует стандартный логгер приложения с указанным уровнем логирования.
// Логи выводятся в формате JSON в стандартный поток вывода (stdout).
//
// Параметры:
//   - logLevel: уровень логирования (slog.Level)
//     Примеры: slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError
//
// Особенности:
//   - Использует JSON формат для структурированного логирования
//   - Заменяет дефолтный логгер приложения
//   - Все последующие вызовы slog.Info(), slog.Error() и т.д. будут использовать этот обработчик
//
// Пример использования:
//
//	setNewDefaultLogger(slog.LevelInfo) // Устанавливает логгер с уровнем Info
func setNewDefaultLogger(logLevel slog.Level) {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})),
	)
}
