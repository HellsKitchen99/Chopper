package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DailyNotesRepository interface {
	CreateNote(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error
	ChangeMood(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error
	ChangeSleepHours(ctx context.Context, userId uuid.UUID, date time.Time, mood float64) error
	ChangeLoad(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error
}
