package domain

import (
	"time"

	"github.com/google/uuid"
)

type SevenDays struct {
	UserId uuid.UUID
	Days   [7]Day
}

type Day struct {
	Date       time.Time
	Mood       int16
	SleepHours float64
	Load       int16
}
