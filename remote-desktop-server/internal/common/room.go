// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

import (
	"time"

	"github.com/google/uuid"
)

// Room представляет модель игровой комнаты в базе данных.
// Поля:
//   - ID: уникальный идентификатор комнаты
//   - Name: название комнаты (4-255 символов)
//   - IsPrivate: флаг приватности комнаты
//   - CreatorID: ID создателя комнаты (uuid)
//   - Password: пароль для приватной комнаты (не возвращается в JSON)
//   - Capacity: максимальное количество игроков
//   - CreatedAt: дата создания комнаты
//   - UpdatedAt: дата обновления (не возвращается в JSON)
//   - DeletedAt: дата удаления (soft delete, не возвращается в JSON)
type Room struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	IsPrivate bool       `json:"is_private"`
	CreatorID uuid.UUID  `json:"creator_id"`
	Password  string     `json:"-"`
	Capacity  uint8      `json:"capacity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// RoomRequest представляет структуру запроса для создания/обновления комнаты.
// Поля с валидацией:
//   - CreatorID: ID создателя комнаты (обязательное)
//   - Name: название комнаты (обязательное, 4-255 символов)
//   - IsPrivate: флаг приватности (обязательное boolean значение)
//   - Password: пароль (обязательное если IsPrivate=true, максимум 255 символов)
type RoomRequest struct {
	CreatorID uuid.UUID `json:"creator_id"`
	Name      string    `validate:"required,min=4,max=255" json:"name"`
	IsPrivate *bool     `validate:"required,boolean" json:"is_private"`
	Password  *string   `validate:"required_if=IsPrivate true,max=255" json:"password"`
}

// RoomResponse представляет упрощенную структуру комнаты для API ответов.
// Поля:
//   - ID: идентификатор комнаты
//   - Name: название комнаты
//   - IsPrivate: флаг приватности
//   - Capacity: вместимость комнаты
//   - PlayerIn: текущее количество игроков в комнате
type RoomResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	IsPrivate *bool  `json:"is_private"`
	Capacity  uint8  `json:"capacity"`
	PlayerIn  int    `json:"player_in"`
}

// RoomSessionResponse представляет полную информацию о комнате для игровой сессии.
// Поля:
//   - ID: идентификатор комнаты
//   - Name: название комнаты
//   - CreatorID: ID создателя комнаты
//   - Password: пароль комнаты (не возвращается в JSON)
//   - IsPrivate: флаг приватности (может быть опущен)
//   - Capacity: вместимость комнаты
//   - Users: список пользователей в комнате (сокращенная информация)
type RoomSessionResponse struct {
	ID        uint64          `json:"id"`
	Name      string          `json:"name"`
	CreatorID uuid.UUID       `json:"creator_id"`
	Password  string          `json:"-"`
	IsPrivate *bool           `json:"is_private,omitempty"`
	Capacity  uint8           `json:"capacity"`
	Users     []*UserResponse `json:"users"`
}
