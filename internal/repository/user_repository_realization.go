package repository

import (
	"chopper/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return err
	}
	return nil
}

func (u *UserRepositoryRealization) CheckUser(ctx context.Context, username string) (domain.User, error) {
	sql := "SELECT id, username, email, password_hash, role, created_at, deleted_at FROM Users WHERE username = $1"
	row := u.pool.QueryRow(ctx, sql, username)
	var user domain.User
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.HashPassword, &user.Role, &user.CreatedAt, &user.DeletedAt); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, ErrNoRow
	} else if err != nil {
		fmt.Println(err)
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepositoryRealization) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	sql := "SELECT id, username, role FROM Users WHERE id = $1 AND username = $2"
	var user domain.UserWhoAmI
	row := u.pool.QueryRow(ctx, sql, id, username)
	if err := row.Scan(&user.Id, &user.Username, &user.Role); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return domain.UserWhoAmI{}, ErrNoRow
	} else if err != nil {
		return domain.UserWhoAmI{}, err
	}
	return user, nil
}
