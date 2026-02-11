package domain

import "github.com/google/uuid"

type UserWhoAmI struct {
	Id       uuid.UUID
	Username string
	Role     Role
}
