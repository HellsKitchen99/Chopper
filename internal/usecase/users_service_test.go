package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repository"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Тест CreateUser - успех
// мок репозитория
type MockUserRepositorySuccess struct {
	wasCalled            bool
	recievedUUID         uuid.UUID
	recievedUsername     string
	recievedEmail        string
	recievedHashPassword string
	recievedRole         domain.Role
}

func (m *MockUserRepositorySuccess) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	m.wasCalled = true
	m.recievedUUID = uuid
	m.recievedUsername = username
	m.recievedEmail = email
	m.recievedHashPassword = hashPassword
	m.recievedRole = role
	return nil
}

func (m *MockUserRepositorySuccess) CheckUser(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockUserRepositorySuccess) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

// мок хэша
type MockPasswordHasherSuccess struct {
	generateWasCalled    bool
	generatePassword     string
	generateHashPassword string
}

func (m *MockPasswordHasherSuccess) GenerateFromPassword(password string) (string, error) {
	m.generateWasCalled = true
	m.generatePassword = password
	m.generateHashPassword = password + "morgan"
	return password + "morgan", nil
}

func (m *MockPasswordHasherSuccess) CompareHashAndPassword(hashPssword string, password string) error {
	return nil
}

// мок генератора UUID
type MockIdGeneratorSuccess struct {
	wasCalled bool
}

func (m *MockIdGeneratorSuccess) NewId() uuid.UUID {
	m.wasCalled = true
	return uuid.UUID{}
}

func TestCreateUserSuccess(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userRegisterFromFront := domain.UserRegisterFromFront{
		Username: "dexter",
		Email:    "dexter@email.com",
		Password: "bay harbour butcher",
	}
	mockUserRepository := &MockUserRepositorySuccess{}
	mockPasswordHasher := &MockPasswordHasherSuccess{}
	mockIdGenerator := &MockIdGeneratorSuccess{}
	service := NewUserService(mockUserRepository, nil, mockPasswordHasher, mockIdGenerator)

	//test
	err := service.CreateUser(ctx, userRegisterFromFront)

	//assert
	if err != nil {
		t.Errorf("ошибки не ожидалось")
	}
	if mockUserRepository.wasCalled != true {
		t.Error("user repository не был вызван")
	}
	if mockUserRepository.recievedUsername != "dexter" {
		t.Errorf("ожидалось имя пользователя - %v", "dexter")
	}
	if mockUserRepository.recievedEmail != "dexter@email.com" {
		t.Errorf("ожидалась почта - %v", "dexter@email.com")
	}
	if mockUserRepository.recievedHashPassword != "bay harbour butcher"+"morgan" {
		t.Errorf("ожидался пароль - %v", "bay harbour butcher"+"morgan")
	}
	if mockUserRepository.recievedRole != domain.RoleUser {
		t.Errorf("ожидалась роль - %v", domain.RoleUser)
	}
	if !mockPasswordHasher.generateWasCalled {
		t.Errorf("password hahser generateFromPassword не был вызван")
	}
	if mockPasswordHasher.generatePassword != userRegisterFromFront.Password {
		t.Errorf("ожидаемый пароль - %v", userRegisterFromFront.Password)
	}
	if mockPasswordHasher.generateHashPassword != userRegisterFromFront.Password+"morgan" {
		t.Errorf("ожидаемый хэш - %v", userRegisterFromFront.Password+"morgan")
	}
}

// Тест CreateUser - провал (длинный пароль)
// Мок хэша
type MockPasswordHasherFailureLongPassword struct {
	wasCalled bool
}

func (m *MockPasswordHasherFailureLongPassword) GenerateFromPassword(password string) (string, error) {
	m.wasCalled = true
	return "", bcrypt.ErrPasswordTooLong
}

func (m *MockPasswordHasherFailureLongPassword) CompareHashAndPassword(hashPssword string, password string) error {
	m.wasCalled = true
	return nil
}

func TestCreateUserFailureLongPassword(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	password := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	userRegisterFromFront := domain.UserRegisterFromFront{
		Username: "dexter",
		Email:    "dexter@email.com",
		Password: password,
	}
	mockUserRepository := &MockUserRepositorySuccess{}
	mockPasswordHasher := &MockPasswordHasherFailureLongPassword{}
	mockIdGenerator := &MockIdGeneratorSuccess{}
	service := NewUserService(mockUserRepository, nil, mockPasswordHasher, mockIdGenerator)

	// test
	err := service.CreateUser(ctx, userRegisterFromFront)

	// assert
	if !errors.Is(err, bcrypt.ErrPasswordTooLong) {
		t.Errorf("ожидалась ошибка - %v", err)
	}
	if mockPasswordHasher.wasCalled != true {
		t.Error("password hasher не был вызван")
	}
	if mockUserRepository.wasCalled != false {
		t.Error("user repository был вызван")
	}
}

// Тест CreateUser - провал (ошибка репозитория)
var MockErrNeedError = errors.New("need error")

// Мок Репозитория
type MockUserRepositoryFailure struct {
}

func (m *MockUserRepositoryFailure) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return MockErrNeedError
}

func (m *MockUserRepositoryFailure) CheckUser(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockUserRepositoryFailure) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

func TestCreateserFailureRepositoryError(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userRegisterFromFront := domain.UserRegisterFromFront{
		Username: "dexter",
		Email:    "dexter@email.com",
		Password: "bay harbour butcher",
	}
	mockUserRepositoryFailure := &MockUserRepositoryFailure{}
	mockPasswordHasherSuccess := &MockPasswordHasherSuccess{}
	mockIdGeneratorSeuccess := &MockIdGeneratorSuccess{}
	service := NewUserService(mockUserRepositoryFailure, nil, mockPasswordHasherSuccess, mockIdGeneratorSeuccess)
	expectedError := MockErrNeedError

	// test
	err := service.CreateUser(ctx, userRegisterFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("ожидалась ошибка - %v", expectedError)
	}
}

// Тест CheckUserInDatabase - успех
var returnedToken string = "need token"

// Мок репозитория
type MockUserRepositorySuccess2 struct {
	wasCalled bool
	username  string
}

func (m *MockUserRepositorySuccess2) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositorySuccess2) CheckUser(ctx context.Context, username string) (domain.User, error) {
	m.wasCalled = true
	m.username = username
	user := domain.User{
		Id:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Username:     "dexter",
		Email:        "dexter@email.com",
		HashPassword: "hashedPassword",
		Role:         domain.RoleUser,
	}
	return user, nil
}

func (m *MockUserRepositorySuccess2) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

// Мок хэша
type MockHashPasswordSuccess2 struct {
	wasCalled    bool
	hashPassword string
	password     string
}

func (m *MockHashPasswordSuccess2) GenerateFromPassword(password string) (string, error) {
	return "", nil
}

func (m *MockHashPasswordSuccess2) CompareHashAndPassword(hashPssword string, password string) error {
	m.wasCalled = true
	m.hashPassword = hashPssword
	m.password = password
	return nil
}

// Мок jwt
type MockJwtServiceSuccess2 struct {
	wasCalled bool
	id        uuid.UUID
	username  string
	email     string
	role      domain.Role
}

func (m *MockJwtServiceSuccess2) GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error) {
	m.wasCalled = true
	m.id = id
	m.username = username
	m.email = email
	m.role = role
	return returnedToken, nil
}

func (m *MockJwtServiceSuccess2) ValidateToken(signedToken string) (*domain.UserClaims, error) {
	return &domain.UserClaims{}, nil
}

func TestCheckUserInDatabaseSuccess(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userLoginFromFront := domain.UserLoginFromFront{
		Username: "dexter",
		Password: "morgan",
	}
	mockUserRepository := &MockUserRepositorySuccess2{}
	mockJwtService := &MockJwtServiceSuccess2{}
	mockHashPassword := &MockHashPasswordSuccess2{}
	service := NewUserService(mockUserRepository, mockJwtService, mockHashPassword, nil)
	expectedId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedUsername := "dexter"
	expectedEmail := "dexter@email.com"
	expectedRole := domain.RoleUser

	// test
	token, err := service.CheckUserInDatabase(ctx, userLoginFromFront)

	// assert
	if err != nil {
		t.Errorf("ошибки не ожидалось")
	}
	if token != returnedToken {
		t.Errorf("ожидался токен - %v", returnedToken)
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не вызван")
	}
	if mockHashPassword.wasCalled != true {
		t.Errorf("hash password не вызван")
	}
	if mockJwtService.wasCalled != true {
		t.Errorf("jwt service не вызван")
	}
	if mockJwtService.id != expectedId {
		t.Errorf("ожидался uuid - %v", expectedId)
	}
	if mockJwtService.username != expectedUsername {
		t.Errorf("ожидалось имя - %v", expectedUsername)
	}
	if mockJwtService.email != expectedEmail {
		t.Errorf("ожидался email - %v", expectedEmail)
	}
	if mockJwtService.role != expectedRole {
		t.Errorf("ожидалась роль - %v", expectedRole)
	}
}

// Тест CheckUserInDatabase - провал (ошибка бд)
var MockErrUserNotExists = errors.New("this error says that user not exists")

// Мок репозитория
type MockUserRepositoryFailureDatabaseError2 struct {
	wasCalled bool
}

func (m *MockUserRepositoryFailureDatabaseError2) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositoryFailureDatabaseError2) CheckUser(ctx context.Context, username string) (domain.User, error) {
	m.wasCalled = true
	return domain.User{}, MockErrUserNotExists
}

func (m *MockUserRepositoryFailureDatabaseError2) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

// Мок хэша
type MockPasswordHashFailureDatabaseError2 struct {
	wasCalled bool
}

func (m *MockPasswordHashFailureDatabaseError2) GenerateFromPassword(password string) (string, error) {
	return "", nil
}

func (m *MockPasswordHashFailureDatabaseError2) CompareHashAndPassword(hashPssword string, password string) error {
	m.wasCalled = true
	return nil
}

// Мок jwt
type MockJwtServiceFailureDatabaseError2 struct {
	wasCalled bool
}

func (m *MockJwtServiceFailureDatabaseError2) GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error) {
	m.wasCalled = true
	return "", nil
}

func (m *MockJwtServiceFailureDatabaseError2) ValidateToken(signedToken string) (*domain.UserClaims, error) {
	return nil, nil
}

func TestCheckUserInDatabaseFailureDatabaseError(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userLoginFromFront := domain.UserLoginFromFront{
		Username: "dexter",
		Password: "the bay harbour butcher",
	}
	mockUserRepository := &MockUserRepositoryFailureDatabaseError2{}
	mockJwtService := &MockJwtServiceFailureDatabaseError2{}
	mockPasswordHash := &MockPasswordHashFailureDatabaseError2{}
	expectedError := MockErrUserNotExists
	service := NewUserService(mockUserRepository, mockJwtService, mockPasswordHash, nil)

	// test
	token, err := service.CheckUserInDatabase(ctx, userLoginFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("ожидалась ошибка - %v", expectedError)
	}
	if token != "" {
		t.Errorf("ожидался пустой токен")
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не вызвался")
	}
	if mockPasswordHash.wasCalled != false {
		t.Errorf("password hash был вызван")
	}
	if mockJwtService.wasCalled != false {
		t.Errorf("jwt service был вызван")
	}
}

// Тест CheckUserInDatabase - провал (ошибка проврки пароля)
var MockErrWrongPassword = errors.New("need error in password")

// Мок репозитория
type MockUserRepositoryFailureWrongPassword3 struct {
	wasCalled bool
}

func (m *MockUserRepositoryFailureWrongPassword3) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositoryFailureWrongPassword3) CheckUser(ctx context.Context, username string) (domain.User, error) {
	m.wasCalled = true
	return domain.User{}, nil
}

func (m *MockUserRepositoryFailureWrongPassword3) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

// Мок хэша
type MockPasswordHashFailureWrongPassword3 struct {
	wasCalled bool
	password  string
}

func (m *MockPasswordHashFailureWrongPassword3) GenerateFromPassword(password string) (string, error) {
	return "", nil
}

func (m *MockPasswordHashFailureWrongPassword3) CompareHashAndPassword(hashPssword string, password string) error {
	m.wasCalled = true
	m.password = password
	return MockErrWrongPassword
}

// Мок jwt
type MockJwtServiceFailureWrongPassword3 struct {
	wasCalled bool
}

func (m *MockJwtServiceFailureWrongPassword3) GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error) {
	m.wasCalled = true
	return "", nil
}

func (m *MockJwtServiceFailureWrongPassword3) ValidateToken(signedToken string) (*domain.UserClaims, error) {
	return nil, nil
}

func TestCheckUserInDatabaseFailureWrongPassword(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userLoginFromFront := domain.UserLoginFromFront{
		Username: "dexter",
		Password: "the bay harbour butcher",
	}
	mockUserRepository := &MockUserRepositoryFailureWrongPassword3{}
	mockJwtService := &MockJwtServiceFailureWrongPassword3{}
	mockPasswordHash := &MockPasswordHashFailureWrongPassword3{}
	service := NewUserService(mockUserRepository, mockJwtService, mockPasswordHash, nil)
	expectedError := ErrWrongPassword

	// test
	token, err := service.CheckUserInDatabase(ctx, userLoginFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("ожидалась ошибка - %v", MockErrWrongPassword)
	}
	if token != "" {
		t.Errorf("ожидался пустой токен")
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не был вызван")
	}
	if mockPasswordHash.wasCalled != true {
		t.Errorf("password hash не был вызван")
	}
	if mockPasswordHash.password != "the bay harbour butcher" {
		t.Errorf("ожидался пароль - %v", "the bay harbour butcher")
	}
	if mockJwtService.wasCalled != false {
		t.Errorf("jwt service был вызван")
	}
}

// Тест CheckUserInDatabase - провал (ошибка генерации токена)
var mockUsername string = "dexter"
var mockPassword string = "morgan"
var mockHashPassword string = "dexter morgan"
var mockEmail string = "dexter@morgan"
var MockErrWhileToken = errors.New("need error while token")

// Мок репозитория
type MockUserRepositoryFailureTokenGeneration4 struct {
	wasCalled bool
	username  string
}

func (m *MockUserRepositoryFailureTokenGeneration4) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositoryFailureTokenGeneration4) CheckUser(ctx context.Context, username string) (domain.User, error) {
	m.wasCalled = true
	m.username = mockUsername
	return domain.User{
		Id:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Username:     mockUsername,
		Email:        mockEmail,
		HashPassword: mockHashPassword,
		Role:         domain.RoleUser,
	}, nil
}

func (m *MockUserRepositoryFailureTokenGeneration4) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	return domain.UserWhoAmI{}, nil
}

// Мок jwt
type MockJwtServiceFailureTokenGeneration4 struct {
	wasCalled bool
	username  string
	email     string
}

func (m *MockJwtServiceFailureTokenGeneration4) GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error) {
	m.wasCalled = true
	m.username = username
	m.email = email
	return "", MockErrWhileToken
}

func (m *MockJwtServiceFailureTokenGeneration4) ValidateToken(signedToken string) (*domain.UserClaims, error) {
	return nil, nil
}

// Мок password hash
type MockPasswordHashFailureTokenGeneration4 struct {
	wasCalled    bool
	hashPassword string
	password     string
}

func (m *MockPasswordHashFailureTokenGeneration4) GenerateFromPassword(password string) (string, error) {
	return "", nil
}

func (m *MockPasswordHashFailureTokenGeneration4) CompareHashAndPassword(hashPssword string, password string) error {
	m.wasCalled = true
	m.hashPassword = hashPssword
	m.password = password
	return nil
}

func TestTestCheckUserInDatabaseFailureTokenError(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userLoginFromFront := domain.UserLoginFromFront{
		Username: mockUsername,
		Password: mockPassword,
	}
	mockUserRepository := &MockUserRepositoryFailureTokenGeneration4{}
	mockJwtService := &MockJwtServiceFailureTokenGeneration4{}
	mockPasswordHash := &MockPasswordHashFailureTokenGeneration4{}
	expectedError := MockErrWhileToken
	service := NewUserService(mockUserRepository, mockJwtService, mockPasswordHash, nil)

	// test
	token, err := service.CheckUserInDatabase(ctx, userLoginFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("ожидалась ошибка - %v", expectedError)
	}
	if token != "" {
		t.Errorf("ожидался пустой токен")
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не был вызван")
	}
	if mockUserRepository.username != mockUsername {
		t.Errorf("ожидалось имя - %v", mockUsername)
	}
	if mockPasswordHash.wasCalled != true {
		t.Errorf("password hash не был вызван")
	}
	if mockPasswordHash.password != mockPassword {
		t.Errorf("ожидался пароль - %v", mockPassword)
	}
	if mockPasswordHash.hashPassword != mockHashPassword {
		t.Errorf("ожидался хэш - %v", mockPasswordHash)
	}
	if mockJwtService.wasCalled != true {
		t.Errorf("jwt service не был вызван")
	}
	if mockJwtService.username != mockUsername {
		t.Errorf("ожидалось имя - %v", mockUsername)
	}
	if mockJwtService.email != mockEmail {
		t.Errorf("ожидался email - %v", mockEmail)
	}
}

// Тест GetIdUsernameRole - Успех
var mockUUID uuid.UUID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var mockUsernameGetIdUsernameRole string = "dexter"
var mockRole domain.Role = domain.RoleUser

// Мок репозитория
type MockUserRepositorySuccess3 struct {
	wasCalled bool
	id        uuid.UUID
	username  string
}

func (m *MockUserRepositorySuccess3) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositorySuccess3) CheckUser(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockUserRepositorySuccess3) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	m.wasCalled = true
	m.id = id
	m.username = username
	return domain.UserWhoAmI{
		Id:       mockUUID,
		Username: mockUsernameGetIdUsernameRole,
		Role:     mockRole,
	}, nil
}

func TestGetIdUsernameRoleSuccess(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id, username := uuid.MustParse("11111111-1111-1111-1111-111111111111"), "dexter"
	mockUserRepository := &MockUserRepositorySuccess3{}
	service := NewUserService(mockUserRepository, nil, nil, nil)

	// test
	user, err := service.GetIdUsernameRole(ctx, id, username)

	// assert
	if err != nil {
		t.Errorf("ошибки не ожидалось")
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не был вызван")
	}
	if mockUserRepository.id != id {
		t.Errorf("ожидался id - %v", mockUUID)
	}
	if mockUserRepository.username != username {
		t.Errorf("ожидался username - %v", mockUsernameGetIdUsernameRole)
	}
	expectedId, expectedUsername, expectedRole := uuid.MustParse("11111111-1111-1111-1111-111111111111"), "dexter", domain.RoleUser
	if user.Id != expectedId {
		t.Errorf("ожидался id - %v", expectedId)
	}
	if user.Username != expectedUsername {
		t.Errorf("ожидался username - %v", expectedUsername)
	}
	if user.Role != expectedRole {
		t.Errorf("ожидался role - %v", expectedRole)
	}
}

// Тест GetIdUsernameRole - провал (ErrNoRows)
var MockErrNoRows = repository.ErrNoRow

// Мок репозитория
type MockUserRepositoryFailureErrNoRows5 struct {
	wasCalled bool
	id        uuid.UUID
	username  string
}

func (m *MockUserRepositoryFailureErrNoRows5) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositoryFailureErrNoRows5) CheckUser(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockUserRepositoryFailureErrNoRows5) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	m.wasCalled = true
	m.id = id
	m.username = username
	return domain.UserWhoAmI{}, MockErrNoRows
}

func TestGetIdUsernameRoleFailureErrNoRows(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockUserRepository := &MockUserRepositoryFailureErrNoRows5{}
	id, username := uuid.MustParse("11111111-1111-1111-1111-111111111111"), "dexter"
	service := NewUserService(mockUserRepository, nil, nil, nil)
	expectedError := ErrUserNotExist

	// test
	_, err := service.GetIdUsernameRole(ctx, id, username)

	// assert
	if !errors.Is(err, expectedError) {
		fmt.Println("ОШИБКА - ", err, expectedError)
		t.Errorf("ожидалась ошибка - %v", expectedError)
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не вызвался")
	}
	if mockUserRepository.id != id {
		t.Errorf("ожидался id - %v", id)
	}
	if mockUserRepository.username != username {
		t.Errorf("ожидался username - %v", username)
	}
}

// Тест GetIdUsernameRole - провал (error)
var MockNeedErr = errors.New("need error")

// Мок репозитория
type MockUserRepositoryFailure6 struct {
	wasCalled bool
	id        uuid.UUID
	username  string
}

func (m *MockUserRepositoryFailure6) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, hashPassword string, role domain.Role) error {
	return nil
}

func (m *MockUserRepositoryFailure6) CheckUser(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (m *MockUserRepositoryFailure6) GetIdUsernameRole(ctx context.Context, id uuid.UUID, username string) (domain.UserWhoAmI, error) {
	m.wasCalled = true
	m.id = id
	m.username = username
	return domain.UserWhoAmI{}, MockNeedErr
}

func TestGetIdUsernameRoleFailureError(t *testing.T) {
	// preparing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id, username := uuid.MustParse("11111111-1111-1111-1111-111111111111"), "dexter"
	mockUserRepository := &MockUserRepositoryFailure6{}
	service := NewUserService(mockUserRepository, nil, nil, nil)
	expectedError := MockNeedErr

	// test
	_, err := service.GetIdUsernameRole(ctx, id, username)

	//assert
	if !errors.Is(err, expectedError) {
		t.Errorf("ожидалась ошибка - %v", expectedError)
	}
	if mockUserRepository.wasCalled != true {
		t.Errorf("user repository не вызвался")
	}
	if mockUserRepository.id != id {
		t.Errorf("ожидался id - %v", id)
	}
	if mockUserRepository.username != username {
		t.Errorf("ожидался username - %v", username)
	}
}
