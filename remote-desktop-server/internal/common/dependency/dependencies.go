// Package dependency управляет зависимостями приложения и внедрением зависимостей.
package dependency

import (
	"log/slog"

	http_handler "github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/handler/http"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/service"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/storage/postgres"
)

// GlobalRepositories содержит все интерфейсы репозиториев, используемые в приложении.
// Служит контейнером для зависимостей слоя доступа к данным.
type GlobalRepositories struct {
	UserRepository  repository.UserRepository  // Репозиторий для операций с пользователями
	ScoreRepository repository.ScoreRepository // Репозиторий для работы с игровыми рекордами
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
	RoomHandler  http_handler.RoomHandler
	ScoreHandler http_handler.ScoreHandler
	UserHandler  http_handler.UserHandler
	AuthHandler  http_handler.AuthHandler
	WSServer     *service.WSServer
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
	db, err := store.NewConnection()
	if err != nil {
		slog.With(op, err.Error())
		panic(err)
	}

	// Инициализация репозиториев
	roomRepo := repository.NewRoomRepository(db)
	scoreRepo := repository.NewScoreRepository(db)
	userRepo := repository.NewUserRepository(db)
	// Инициализация сервисов
	roomService := service.NewRoomService(roomRepo)
	scoreService := service.NewScoreService(scoreRepo, userRepo)
	userService := service.NewUserService(userRepo, scoreRepo)
	authService := service.NewAuthService(userRepo)
	// Создание обработчиков
	roomHandler := http_handler.NewRoomHandler(*roomService)
	scoreHandler := http_handler.NewScoreHandler(*scoreService)
	userHandler := http_handler.NewUserHandler(*userService)
	authHandler := http_handler.NewAuthHandler(*authService)

	return &AppDependencies{
		RoomHandler:  *roomHandler,
		ScoreHandler: *scoreHandler,
		UserHandler:  *userHandler,
		AuthHandler:  *authHandler,
		WSServer: service.NewWsServer(
			service.NewScoreService(scoreRepo, userRepo),
		),
		GlobalRepositories: GlobalRepositories{
			UserRepository:  userRepo,
			ScoreRepository: scoreRepo,
		},
	}
}
