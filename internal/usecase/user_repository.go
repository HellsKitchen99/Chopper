package usecase

import (
	"chopper/internal/domain"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error
}
