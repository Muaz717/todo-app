package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Muaz717/todo-app/internal/app/storage"
	"github.com/Muaz717/todo-app/internal/domain/models"

	"github.com/jackc/pgx"
)

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "postgres.SaveUser"

	query := `INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`

	row := s.db.QueryRow(ctx, query, email, passHash)

	var userId int64

	err := row.Scan(&userId)
	if err != nil {
		if pgErr, ok := err.(*pgx.PgError); ok {
			return 0, fmt.Errorf("%s: SQL Error: %s, Detail: %s, Where: %s", op, pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.User"

	query := `SELECT id, email, pass_hash FROM users WHERE email=$1`

	row := s.db.QueryRow(ctx, query, email)

	var user models.User

	err := row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
