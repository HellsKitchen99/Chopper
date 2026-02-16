package usecase

import (
	"chopper/internal/domain"

	"github.com/google/uuid"
)

type JwtGenerator interface {
	GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error)
	ValidateToken(signedToken string) (*domain.UserClaims, error)
}
