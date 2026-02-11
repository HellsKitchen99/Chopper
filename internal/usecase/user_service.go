package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repository"
	"chopper/internal/security"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository UserRepository
	jwtService     *security.Jwt
}

func NewUserService(userRepository UserRepository, jwtService *security.Jwt) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwtService:     jwtService,
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
	if err := u.userRepository.CreateUser(ctx, uuid, username, email, string(passwordHash), defaultRole); err != nil && errors.Is(err, repository.ErrUniqueViolation) {
		return ErrUserExists
	} else if err != nil {
		return err
	}
	return nil
}

func (u *UserService) CheckUserInDatabase(ctx context.Context, userLoginFromFront domain.UserLoginFromFront) (string, error) {
	username := userLoginFromFront.Username
	user, err := u.userRepository.CheckUser(ctx, username)
	if err != nil && errors.Is(err, repository.ErrNoRow) {
		return "", ErrUserNotExist
	} else if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(userLoginFromFront.Password)); err != nil {
		return "", ErrWrongPassword
	}
	token, err := u.jwtService.GenerateToken(user.Id, user.Username, user.Email, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	user, err := u.userRepository.GetIdUsernameRole(ctx, id, username)
	if err != nil && errors.Is(err, repository.ErrNoRow) {
		return user, ErrUserNotExist
	} else if err != nil {
		return user, err
	}
	return user, nil
}
