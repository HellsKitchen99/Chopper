package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	return &PasswordHasher{
		cost: cost,
	}
}

func (p *PasswordHasher) GenerateFromPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func (p *PasswordHasher) CompareHashAndPassword(hashPssword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPssword), []byte(password)); err != nil {
		return err
	}
	return nil
}
