package repository

import (
	"chopper/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateUserSuccess(t *testing.T) {
	// preparing
	ctx, id := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	username, email, hashPassword := "dexter", "thebayharbourbutcher", "hash_password"
	role := domain.RoleUser
	userRepo := NewUserRepositoryRealization(testPool)

	// test
	err := userRepo.CreateUser(ctx, id, username, email, hashPassword, role)

	// assert
	if err != nil {
		t.Errorf("expected error was nil")
	}
	sql := "SELECT 1 FROM Users WHERE id = $1"
	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	row := userRepo.pool.QueryRow(c, sql, id)
	var result int
	if err := row.Scan(&result); err != nil {
		t.Errorf("error while trying to scan result")
	}
	if result != 1 {
		t.Errorf("expected result was 1")
	}
}
