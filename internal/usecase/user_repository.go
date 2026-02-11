package usecase

import (
	"chopper/internal/domain"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error
	CheckUser(ctx context.Context, username string) (domain.User, error)
	GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error)
}
