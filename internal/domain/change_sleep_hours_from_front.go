package domain

import "time"

type ChangeSleepHoursFromFront struct {
	Date       time.Time `json:"date"`
	SleepHours float64   `json:"sleep_hours"`
}
