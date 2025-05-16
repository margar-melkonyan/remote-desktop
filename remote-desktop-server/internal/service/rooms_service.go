// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/config"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RoomService предоставляет методы для управления игровыми комнатами.
type RoomService struct {
	repo repository.RoomRepository
}

const op = "room_serivce"

// NewRoomService создаёт новый экземпляр RoomService с указанным репозиторием.
func NewRoomService(repoRoom repository.RoomRepository) *RoomService {
	return &RoomService{
		repo: repoRoom,
	}
}

// GetAll возвращает список всех доступных комнат, в которых есть свободные места.
func (service *RoomService) GetAll(ctx context.Context, ws *WSServer) []*common.RoomResponse {
	rooms, err := service.repo.FindAll(ctx)
	var roomsResponse []*common.RoomResponse
	roomsResponse = make([]*common.RoomResponse, 0)

	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	for _, room := range rooms {
		playerIn := 0
		roomInfo := ws.Rooms[room.ID]
		if roomInfo != nil {
			playerIn = len(roomInfo.Users)
		}
		if playerIn != 2 {
			roomsResponse = append(roomsResponse, &common.RoomResponse{
				ID:        room.ID,
				Name:      room.Name,
				Capacity:  room.Capacity,
				IsPrivate: &room.IsPrivate,
				PlayerIn:  playerIn,
			})
		}
	}
	if err != nil {
		slog.Error(
			fmt.Sprintf("[%v:GetAll]", op),
			slog.String("error", err.Error()),
		)
		return []*common.RoomResponse{}
	}
	return roomsResponse
}

// isUserInRoom проверяет, присутствует ли пользователь в списке подключённых к комнате пользователей.
func isUserInRoom(currentUser *common.User, users []*ConnectedUser) bool {
	for _, user := range users {
		if user.ID == currentUser.ID {
			return true
		}
	}
	return false
}

// GetAllMy возвращает комнаты, созданные или занятые текущим пользователем.
func (service *RoomService) GetAllMy(ctx context.Context, ws *WSServer) []*common.RoomResponse {
	currentUser, ok := ctx.Value(common.USER).(*common.User)
	if !ok {
		return nil
	}
	rooms, err := service.repo.FindAll(ctx)
	var roomsResponse []*common.RoomResponse
	roomsResponse = make([]*common.RoomResponse, 0)
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	for _, room := range rooms {
		playerIn := 0
		roomInfo := ws.Rooms[room.ID]
		if roomInfo != nil {
			playerIn = len(roomInfo.Users)
		}
		if currentUser.ID == room.CreatorID || (roomInfo != nil && isUserInRoom(currentUser, roomInfo.Users)) {
			roomsResponse = append(roomsResponse, &common.RoomResponse{
				ID:        room.ID,
				Name:      room.Name,
				Capacity:  room.Capacity,
				IsPrivate: &room.IsPrivate,
				PlayerIn:  playerIn,
			})
		}
	}
	if err != nil {
		slog.Error(
			fmt.Sprintf("[%v:GetAll]", op),
			slog.String("error", err.Error()),
		)
		return []*common.RoomResponse{}
	}
	return roomsResponse
}

// GetById возвращает информацию о комнате по её идентификатору, включая подключённых пользователей.
func (service *RoomService) GetById(ctx context.Context, id uint64, ws *WSServer) (*common.RoomSessionResponse, error) {
	room, err := service.repo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find room with ID %d: %w", id, err)
	}
	if room == nil {
		return nil, fmt.Errorf("room with ID %d not found", id)
	}
	users := make([]*common.UserResponse, 0)
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	roomSession := ws.Rooms[room.ID]
	if roomSession != nil {
		for _, user := range roomSession.Users {
			users = append(users, &common.UserResponse{
				ID:     user.ID,
				Name:   user.Name,
				Symbol: user.Symbol,
			})
		}
	}
	resp := &common.RoomSessionResponse{
		ID:        room.ID,
		Name:      room.Name,
		CreatorID: room.CreatorID,
		Password:  room.Password,
		IsPrivate: &room.IsPrivate,
		Capacity:  room.Capacity,
		Users:     users,
	}
	return resp, nil
}

// Create создаёт новую игровую комнату. Если установлен пароль, он хэшируется с помощью bcrypt.
func (service *RoomService) Create(ctx context.Context, form common.RoomRequest) error {
	if *form.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(*form.Password), config.ServerConfig.BcryptPower)
		if err != nil {
			return err
		}
		hashedPassword := string(password)
		form.Password = &hashedPassword
	}
	if *form.Password == "" {
		*form.IsPrivate = false
	}
	user, ok := ctx.Value(common.USER).(*common.User)
	if !ok {
		return errors.New("userId is not correct")
	}
	room := common.Room{
		CreatorID: user.ID,
		Name:      form.Name,
		Password:  *form.Password,
		IsPrivate: *form.IsPrivate,
	}
	return service.repo.Create(ctx, room)
}

// DeleteById удаляет комнату по её идентификатору.
func (service *RoomService) DeleteById(ctx context.Context, id uint64) error {
	return service.repo.DeleteById(ctx, id)
}
