package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Username  string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	DeletedAt time.Time
}
