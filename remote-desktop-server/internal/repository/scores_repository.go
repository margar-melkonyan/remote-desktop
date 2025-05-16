// Package repository предоставляет реализации репозиториев для работы с данными приложения.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
)

const TABLE_NAME = "scores" // Название таблицы в базе данных

// ScoreRepo реализует ScoreRepository для работы с PostgreSQL
type ScoreRepo struct {
	db *sql.DB
}

// ScoreRepository определяет контракт для работы с хранилищем результатов игр
type ScoreRepository interface {
	// Create создает новую запись о результате игры
	Create(ctx context.Context, score *common.Score) error

	// FindAllByUser возвращает последние 50 результатов игр для указанного пользователя
	FindAllByUser(ctx context.Context, user *common.User) ([]*common.Score, error)

	// GetWonScore возвращает количество побед указанного пользователя
	GetWonScore(ctx context.Context, user *common.User) (uint, error)
}

// NewScoreRepository создает новый экземпляр ScoreRepository
func NewScoreRepository(db *sql.DB) ScoreRepository {
	return &ScoreRepo{
		db: db,
	}
}

// Create сохраняет результат игры в базе данных
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - score: указатель на структуру с данными результата
//
// Возвращает:
//   - error: ошибка, если не удалось создать запись
//
// Особенности:
//   - Сохраняет nickname, user_id и флаг победы (is_won)
//   - Проверяет количество затронутых строк (rowsAffected)
//   - Возвращает ошибку "room was not created" если не была создана запись
func (repo ScoreRepo) Create(ctx context.Context, score *common.Score) error {
	query := "INSERT INTO scores (name, user_id, is_won) VALUES ($1, $2, $3)"
	result, err := repo.db.ExecContext(ctx, query, score.Nickname, score.UserID, score.IsWon)
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

// FindAllByUser возвращает историю результатов для пользователя
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - user: указатель на структуру пользователя
//
// Возвращает:
//   - []*common.Score: слайс указателей на результаты (максимум 50)
//   - error: ошибка выполнения запроса
//
// Особенности:
//   - Возвращает только активные записи (deleted_at IS NULL)
//   - Сортирует по дате создания (новые сначала)
//   - Если результатов нет, возвращает пустой слайс (не nil)
//   - Ограничивает выборку 50 последними записями
func (repo ScoreRepo) FindAllByUser(ctx context.Context, user *common.User) ([]*common.Score, error) {
	var scores []*common.Score
	query := fmt.Sprintf(
		"SELECT id, name, user_id, is_won, created_at FROM %v WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 50",
		TABLE_NAME,
	)
	rows, err := repo.db.QueryContext(ctx, query, user.ID)
	defer func() {
		rows.Close()
	}()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var score common.Score
		err := rows.Scan(
			&score.ID,
			&score.Nickname,
			&score.UserID,
			&score.IsWon,
			&score.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, &score)
	}
	if scores == nil {
		scores = []*common.Score{}
	}
	return scores, nil
}

// GetWonScore возвращает количество побед пользователя
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - user: указатель на структуру пользователя
//
// Возвращает:
//   - uint: количество побед
//   - error: ошибка выполнения запроса
//
// Особенности:
//   - Считает только записи с is_won = true
//   - Не учитывает удаленные записи (deleted_at IS NULL)
//   - Возвращает 0 в случае ошибки
func (repo ScoreRepo) GetWonScore(ctx context.Context, user *common.User) (uint, error) {
	var currentScore *uint
	query := fmt.Sprintf(
		"SELECT COUNT(*) as current_score FROM %v WHERE user_id = $1 AND is_won = 1",
		TABLE_NAME,
	)
	row := repo.db.QueryRowContext(ctx, query, user.ID)
	err := row.Scan(
		&currentScore,
	)
	if err != nil {
		return 0, err
	}
	return *currentScore, nil
}
