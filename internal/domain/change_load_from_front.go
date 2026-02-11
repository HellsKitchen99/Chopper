package domain

import "time"

type ChangeLoadFromFront struct {
	Date time.Time `json:"date"`
	Load int16     `json:"load"`
}
