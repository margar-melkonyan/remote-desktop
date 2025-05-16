// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"errors"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
)

// ScoreService предоставляет методы для работы со счетами пользователей.
type ScoreService struct {
	scoreRepo repository.ScoreRepository
	userRepo  repository.UserRepository
}

// NewScoreService создаёт новый экземпляр ScoreService.
func NewScoreService(
	scoreRepo repository.ScoreRepository,
	userRepo repository.UserRepository,
) *ScoreService {
	return &ScoreService{
		scoreRepo: scoreRepo,
		userRepo:  userRepo,
	}
}

// GetCurrentUserScores возвращает список результатов текущего пользователя на основе email в контексте.
func (service *ScoreService) GetCurrentUserScores(ctx context.Context) ([]*common.Score, error) {
	email, ok := ctx.Value(common.USER_MAIL).(string)
	if !ok {
		return nil, errors.New("user email is not valid")
	}
	user, err := service.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return service.scoreRepo.FindAllByUser(ctx, user)
}
