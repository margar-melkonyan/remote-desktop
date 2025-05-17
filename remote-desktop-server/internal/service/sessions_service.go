// Package service реализует бизнес-логику приложения.
package service

import (
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
)

// SessionService предоставляет сервис для работы c подключениями к удаленным серверам.
type SessionService struct {
	repo repository.UserRepository
}

// NewSessionService создает новый экземпляр SessionService.
//
// Параметры:
//   - repo: репозиторий для работы с подключениями
//
// Возвращает:
//   - *SessionService: указатель на созданный сервис
func NewSessionService(repo repository.UserRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}
