// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

import (
	"time"
)

// Score представляет модель игрового результата (счета) пользователя.
//
// Поля:
//   - ID: уникальный идентификатор записи (не возвращается в JSON)
//   - UserID: идентификатор пользователя в формате UUID (обязательное поле, не возвращается в JSON)
//   - IsWon: флаг победы (1 - победа, 0 - поражение, обязательное поле)
//   - Nickname: никнейм игрока (отображается в таблице результатов)
//   - CreatedAt: дата создания записи (может быть опущена в JSON)
//
// Валидация:
//   - UserID должен быть валидным UUID
//   - IsWon должен быть булевым значением (0 или 1)
type Score struct {
	ID        uint64    `json:"-"`
	UserID    string    `json:"-" validate:"required,uuid"`
	IsWon     float64   `json:"is_won" validate:"required,boolean"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
