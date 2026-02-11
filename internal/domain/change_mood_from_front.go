package domain

import "time"

type ChangeMoodFromFront struct {
	Date time.Time `json:"date"`
	Mood int16     `json:"mood"`
}
