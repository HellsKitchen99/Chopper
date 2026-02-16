package security

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) NewId() uuid.UUID {
	return uuid.New()
}
