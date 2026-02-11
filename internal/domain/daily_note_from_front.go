package domain

type DailyNoteFromFront struct {
	Mood       int16   `json:"mood"`
	SleepHours float64 `json:"sleep_hours"`
	Load       int16   `json:"load"`
}
