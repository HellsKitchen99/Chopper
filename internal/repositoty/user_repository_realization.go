package repositoty

import (
	"chopper/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryRealization struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryRealization(pool *pgxpool.Pool) *UserRepositoryRealization {
	return &UserRepositoryRealization{
		pool: pool,
	}
}

func (u *UserRepositoryRealization) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	sql := "INSERT INTO Users (id, username, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := u.pool.Exec(ctx, sql, uuid, username, email, hashPassword, role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
	}
	return nil
}
