// Package repository предоставляет реализации репозиториев для работы с данными приложения.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
)

// guacamoleRepo реализует GuacamoleRepository для работы с PostgreSQL
type guacamoleRepo struct {
	db *sql.DB
}

// GuacamoleRepository определяет контракт для работы с хранилищем Apache
type GuacamoleRepository interface {
	CreateEntity(ctx context.Context, username string) (uint64, error)
	CreateUserAndPermissions(ctx context.Context, form common.GuacamoleUser) error
	AddPermissionToUser(ctx context.Context, id int, permissions []string) error
}

// NewUserRepository создает новый экземпляр GuacamoleRepository
func NewGuacamoleRepository(db *sql.DB) GuacamoleRepository {
	return &guacamoleRepo{
		db: db,
	}
}

func (repo *guacamoleRepo) CreateEntity(ctx context.Context, username string) (uint64, error) {
	query := "INSERT INTO guacamole_entity (name, type) VALUES ($1, $2) RETURNING entity_id"
	var id uint64
	err := repo.db.QueryRowContext(
		ctx,
		query,
		username,
		"USER",
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *guacamoleRepo) CreateUserAndPermissions(ctx context.Context, form common.GuacamoleUser) error {
	tx, _ := repo.db.BeginTx(ctx, nil)
	defer tx.Commit()
	query := `
		INSERT INTO guacamole_user (
			entity_id, password_hash, password_salt, password_date, disabled, expired
		) VALUES ($1, decode($2, 'hex'), decode($3, 'hex'), CURRENT_TIMESTAMP, false, false)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		form.ID,
		form.PasswordHex,
		form.SaultHex,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("guacamole account was not created")
	}

	query = `
		INSERT INTO guacamole_system_permission (
			entity_id, permission
		) VALUES ($1, $2)
	`

	result, err = repo.db.ExecContext(
		ctx,
		query,
		form.ID,
		form.Permissions[0],
	)
	if err != nil {
		return err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("permissions for guacamole account was not created")
	}
	return nil
}

func (repo *guacamoleRepo) AddPermissionToUser(ctx context.Context, id int, permissions []string) error {
	return nil
}
