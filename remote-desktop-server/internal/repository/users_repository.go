// Package repository предоставляет реализации репозиториев для работы с данными приложения.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/margar-melkonyan/remote-desktop.git/internal/common"
)

// userRepo реализует UserRepository для работы с PostgreSQL
type userRepo struct {
	db *sql.DB
}

// UserRepository определяет контракт для работы с хранилищем пользователей
type UserRepository interface {
	// FindByEmail находит пользователя по email
	FindByEmail(ctx context.Context, email string) (*common.User, error)

	// Create создает нового пользователя в системе
	Create(ctx context.Context, form common.AuthSignUpRequest) error
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// FindByEmail ищет пользователя по email адресу
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - email: email адрес пользователя
//
// Возвращает:
//   - *common.User: найденный пользователь
//   - error: ошибка если пользователь не найден или произошла ошибка запроса
//
// Особенности:
//   - Возвращает только активных пользователей (deleted_at IS NULL)
//   - Включает в результат: ID, имя, email, хэш пароля и дату создания
func (repo *userRepo) FindByEmail(ctx context.Context, email string) (*common.User, error) {
	var user common.User
	query := "SELECT id, name, email, password, created_at FROM users WHERE email = $1 AND deleted_at IS NULL"
	row := repo.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create регистрирует нового пользователя в системе
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - form: данные для регистрации (AuthSignUpRequest)
//
// Возвращает:
//   - error: ошибка если не удалось создать пользователя
//
// Особенности:
//   - Сохраняет имя, email и хэш пароля
//   - Проверяет количество затронутых строк (rowsAffected)
//   - Возвращает ошибку "room was not created" если запись не была создана
//     (Примечание: возможно стоит изменить текст ошибки на более подходящий)
func (repo *userRepo) Create(ctx context.Context, form common.AuthSignUpRequest) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	result, err := repo.db.ExecContext(
		ctx,
		query,
		&form.Name,
		&form.Email,
		&form.Password,
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
