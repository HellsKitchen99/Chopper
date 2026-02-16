package usecase

import "github.com/google/uuid"

type UUIDGenerator interface {
	NewId() uuid.UUID
}
