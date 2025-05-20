// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"errors"

	"github.com/margar-melkonyan/remote-desktop.git/internal/common"
	"github.com/margar-melkonyan/remote-desktop.git/internal/repository"
)

// UserService предоставляет методы для получения информации о пользователе.
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetCurrentUser возвращает информацию о текущем пользователе, включая счёт побед.
func (service *UserService) GetCurrentUser(ctx context.Context) (*common.UserResponse, error) {
	email, ok := ctx.Value(common.USER_MAIL).(string)
	if !ok {
		return nil, errors.New("user email is not valid")
	}
	user, err := service.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	userResponse := &common.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
	}
	return userResponse, nil
}
