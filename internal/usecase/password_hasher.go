package usecase

type PasswordHasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashPssword string, password string) error
}
