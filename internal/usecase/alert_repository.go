package usecase

import (
	"chopper/internal/domain"
	"context"

	"github.com/google/uuid"
)

type AlertRepository interface {
	GetLastSevenDays(ctx context.Context, userId uuid.UUID) ([]domain.Day, error)
}
