// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

import (
	"time"

	"github.com/google/uuid"
)

// Константы для ключей контекста:
const USER_MAIL = "user_mail"   // Ключ для хранения email пользователя в контексте
const GUAC_TOKEN = "guac_token" // Тоекн для запросов к Apache Guacamole
const USER = "user"             // Ключ для хранения данных пользователя в контексте

// AuthSignInRequest представляет структуру запроса для входа пользователя.
// Поля:
//   - Email: электронная почта (обязательное, 4-255 символов)
//   - Password: пароль (обязательное, 8-64 символов)
type AuthSignInRequest struct {
	Email    string `json:"email" validate:"required,min=4,max=255"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// AuthSignUpRequest представляет структуру запроса для регистрации пользователя.
// Поля:
//   - Name: имя пользователя (обязательное, 4-255 символов)
//   - Email: электронная почта (обязательное, валидный email, 8-255 символов)
//   - Password: пароль (обязательное, 8-255 символов)
//   - PasswordConfirmation: подтверждение пароля (должно совпадать с Password)
type AuthSignUpRequest struct {
	Name                 string `json:"name" validate:"required,min=4,max=255"`
	Email                string `json:"email" validate:"required,email,min=8,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=255"`
}

// User представляет модель пользователя в системе.
// Поля:
//   - ID: уникальный идентификатор (не возвращается в JSON)
//   - Name: имя пользователя
//   - Email: электронная почта
//   - Password: хэш пароля (не возвращается в JSON)
//   - CreatedAt: дата создания
//   - UpdatedAt: дата обновления (не возвращается в JSON)
//   - DeletedAt: дата удаления (soft delete, не возвращается в JSON)
type User struct {
	ID        uuid.UUID  `json:"-"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// UserResponse представляет структуру ответа с данными пользователя.
// Используется для безопасного возврата данных пользователя (без чувствительной информации).
// Поля:
//   - ID: уникальный идентификатор
//   - Name: имя пользователя
//   - Email: электронная почта (может быть опущена)
//   - CreatedAt: дата создания аккаунта (может быть опущена)
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
