package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository UserRepository
	jwtService     JwtGenerator
	passwordHasher PasswordHasher
	uuidGenerator  UUIDGenerator
}

func NewUserService(userRepository UserRepository, jwtService JwtGenerator, passwordHasher PasswordHasher, uuidGenerator UUIDGenerator) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwtService:     jwtService,
		passwordHasher: passwordHasher,
		uuidGenerator:  uuidGenerator,
	}
}

func (u *UserService) CreateUser(ctx context.Context, userRegisterFromFront domain.UserRegisterFromFront) error {
	defaultRole := domain.RoleUser
	username := userRegisterFromFront.Username
	email := userRegisterFromFront.Email
	passwordHash, err := u.passwordHasher.GenerateFromPassword(userRegisterFromFront.Password)
	if err != nil {
		return err
	}
	uuid := u.uuidGenerator.NewId()
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
	if err := u.passwordHasher.CompareHashAndPassword(user.HashPassword, userLoginFromFront.Password); err != nil {
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
