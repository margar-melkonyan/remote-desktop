// Package dependency управляет зависимостями приложения и внедрением зависимостей.
package dependency

import (
	"fmt"
	"log/slog"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/config"
	http_handler "github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/handler/http"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/storage/postgres"
)

// GlobalRepositories содержит все интерфейсы репозиториев, используемые в приложении.
// Служит контейнером для зависимостей слоя доступа к данным.
type GlobalRepositories struct {
	UserRepository repository.UserRepository // Репозиторий для операций с пользователями
}

// AppDependencies содержит все зависимости приложения:
//   - Обработчики HTTP запросов
//   - WebSocket сервер
//   - Глобальные репозитории
//
// Используется для:
//   - Инициализации всех компонентов приложения
//   - Внедрения зависимостей между слоями
//   - Предоставления единой точки доступа к сервисам
type AppDependencies struct {
	UserHandler    http_handler.UserHandler
	AuthHandler    http_handler.AuthHandler
	SessionHandler http_handler.SessionHandler
	GlobalRepositories
}

// NewAppDependencies создает и инициализирует все зависимости приложения.
//
// Выполняет:
//  1. Подключение к базе данных
//  2. Инициализацию репозиториев
//  3. Создание сервисов
//  4. Инициализацию обработчиков
//  5. Настройку WebSocket сервера
//
// Возвращает:
// - *AppDependencies: указатель на инициализированные зависимости
//
// При ошибках подключения к БД завершает работу с panic.
func NewAppDependencies() *AppDependencies {
	const op = "config.NewAppDependencides"
	store := postgres.Storage{
		ConnectionDriver: "postgres",
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.ServerConfig.DbConfig[0].Host,
		config.ServerConfig.DbConfig[0].Port,
		config.ServerConfig.DbConfig[0].Username,
		config.ServerConfig.DbConfig[0].Password,
		config.ServerConfig.DbConfig[0].Name,
		config.ServerConfig.DbConfig[0].SSLMode,
	)
	db, err := store.NewConnection(dsn)
	if err != nil {
		slog.With(op, err.Error())
		panic(err)
	}

	dsnGuacamole := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.ServerConfig.DbConfig[1].Host,
		config.ServerConfig.DbConfig[1].Port,
		config.ServerConfig.DbConfig[1].Username,
		config.ServerConfig.DbConfig[1].Password,
		config.ServerConfig.DbConfig[1].Name,
		config.ServerConfig.DbConfig[1].SSLMode,
	)
	dbGuac, err := store.NewConnection(dsnGuacamole) // dbGUAC
	if err != nil {
		slog.With(op, err.Error())
		panic(err)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	guacRepo := repository.NewGuacamoleRepository(dbGuac)
	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, guacRepo)
	sessionService := service.NewSessionService(userRepo)
	// Создание обработчиков
	userHandler := http_handler.NewUserHandler(*userService)
	authHandler := http_handler.NewAuthHandler(*authService)
	sessionHandler := http_handler.NewSessionHandler(*sessionService)

	return &AppDependencies{
		UserHandler:    *userHandler,
		AuthHandler:    *authHandler,
		SessionHandler: *sessionHandler,
		GlobalRepositories: GlobalRepositories{
			UserRepository: userRepo,
		},
	}
}
