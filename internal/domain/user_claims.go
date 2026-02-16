package domain

import (
	"github.com/google/uuid"
)

type UserClaims struct {
	Id       uuid.UUID
	Username string
	Email    string
	Role     Role
}
