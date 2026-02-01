package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repositoty"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, userRegisterFromFront domain.UserRegisterFromFront) error {
	defaultRole := domain.RoleUser
	username := userRegisterFromFront.Username
	email := userRegisterFromFront.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userRegisterFromFront.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error while trying to hash password: %v", err)
	}
	uuid := uuid.New()
	if err := u.userRepository.CreateUser(ctx, uuid, username, email, string(passwordHash), defaultRole); err != nil && errors.Is(err, repositoty.ErrUniqueViolation) {
		return ErrUserExists
	} else if err != nil {
		return err
	}
	return nil
}
