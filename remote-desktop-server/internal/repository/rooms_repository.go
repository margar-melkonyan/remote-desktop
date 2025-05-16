// Package repository предоставляет реализации репозиториев для работы с данными приложения.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
)

// RoomRepo реализует RoomRepository для работы с PostgreSQL
type RoomRepo struct {
	db *sql.DB
}

// RoomRepository определяет контракт для работы с хранилищем комнат
type RoomRepository interface {
	// FindAll возвращает все активные комнаты (не помеченные как удаленные)
	FindAll(ctx context.Context) ([]*common.Room, error)

	// FindById находит комнату по идентификатору
	FindById(ctx context.Context, id uint64) (*common.Room, error)

	// Create создает новую комнату в базе данных
	Create(ctx context.Context, room common.Room) error

	// DeleteById помечает комнату как удаленную (soft delete)
	DeleteById(ctx context.Context, id uint64) error
}

// NewRoomRepository создает новый экземпляр RoomRepository
func NewRoomRepository(db *sql.DB) RoomRepository {
	return &RoomRepo{
		db: db,
	}
}

// FindAll возвращает список всех активных комнат
//
// Параметры:
//   - ctx: контекст выполнения запроса
//
// Возвращает:
//   - []*common.Room: слайс указателей на комнаты
//   - error: ошибка, если не удалось выполнить запрос
//
// Особенности:
//   - Возвращает только комнаты, где deleted_at IS NULL
//   - Если комнат нет, возвращает пустой слайс (не nil)
func (repo *RoomRepo) FindAll(ctx context.Context) ([]*common.Room, error) {
	var rooms []*common.Room
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM rooms WHERE deleted_at IS NULL")
	defer func() {
		rows.Close()
	}()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var room common.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.IsPrivate,
			&room.Password,
			&room.CreatorID,
			&room.Capacity,
			&room.CreatedAt,
			&room.UpdatedAt,
			&room.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	if rooms == nil {
		rooms = []*common.Room{}
	}
	return rooms, nil
}

// FindById находит комнату по идентификатору
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - id: идентификатор комнаты
//
// Возвращает:
//   - *common.Room: найденная комната
//   - error: ошибка, если комната не найдена или произошла ошибка запроса
//
// Особенности:
//   - Возвращает только активные комнаты (deleted_at IS NULL)
//   - Не выбирает поля updated_at и deleted_at
func (repo *RoomRepo) FindById(ctx context.Context, id uint64) (*common.Room, error) {
	var room common.Room
	query := "SELECT id, name, is_private, password, creator_id, capacity, created_at FROM rooms WHERE id = $1 AND deleted_at IS NULL"
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.Name,
		&room.IsPrivate,
		&room.Password,
		&room.CreatorID,
		&room.Capacity,
		&room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Create создает новую комнату в базе данных
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - room: данные комнаты для создания
//
// Возвращает:
//   - error: ошибка, если не удалось создать комнату
//
// Особенности:
//   - Обязательные поля: name, is_private, creator_id
//   - Поле password может быть пустым для публичных комнат
//   - Проверяет количество затронутых строк (rowsAffected)
func (repo *RoomRepo) Create(ctx context.Context, room common.Room) error {
	query := "INSERT INTO rooms (name, is_private, creator_id, password) VALUES ($1, $2, $3, $4)"
	result, err := repo.db.ExecContext(
		ctx,
		query,
		room.Name,
		room.IsPrivate,
		room.CreatorID,
		room.Password,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("room was not created")
	}
	return nil
}

// DeleteById выполняет мягкое удаление комнаты (soft delete)
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - id: идентификатор комнаты для удаления
//
// Возвращает:
//   - error: ошибка, если не удалось выполнить удаление
//
// Особенности:
//   - Устанавливает deleted_at в текущее время
//   - Не удаляет запись физически
//   - Проверяет количество затронутых строк (rowsAffected)
func (repo *RoomRepo) DeleteById(ctx context.Context, id uint64) error {
	result, err := repo.db.ExecContext(
		ctx,
		"UPDATE rooms SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("room was not deleted")
	}
	return nil
}
