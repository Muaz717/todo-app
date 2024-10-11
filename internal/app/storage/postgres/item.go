package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Muaz717/todo-app/internal/app/storage"
	"github.com/Muaz717/todo-app/internal/domain/models"
	pgx "github.com/jackc/pgx"
	pgx5 "github.com/jackc/pgx/v5"
)

func (s *Storage) SaveItem(
	ctx context.Context,
	userId int64,
	title string,
	description string,
) (int64, error) {
	const op = "postgres.SaveItem"

	query := `INSERT INTO items(title, description, user_id) VALUES($1, $2, $3) RETURNING id`

	row := s.db.QueryRow(ctx, query, title, description, userId)

	var itemId int64
	err := row.Scan(&itemId)
	if err != nil {
		if pgErr, ok := err.(*pgx.PgError); ok {
			return 0, fmt.Errorf("%s: SQL Error: %s, Detail: %s, Where: %s", op, pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return itemId, err
}

func (s *Storage) AllItems(ctx context.Context, userId int64) ([]models.Item, error) {
	const op = "postgres.AllItems"

	query := `SELECT id, title, description FROM items WHERE user_id = $1`

	rows, err := s.db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgx5.CollectRows(rows, pgx5.RowToStructByName[models.Item])
}
