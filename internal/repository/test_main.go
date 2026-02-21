package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var testPool *pgxpool.Pool
var migrationsPath = "file://../../migrations" // internal//repository

func TestMain(m *testing.M) {
	testDataBaseConfig, err := TestDbConfigLoad()
	if err != nil {
		logrus.Fatalf("error while trying to load test database config: %v", err)
	}
	testUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", testDataBaseConfig.User, testDataBaseConfig.Password, testDataBaseConfig.Host, testDataBaseConfig.Port, testDataBaseConfig.Name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	testPool, err = pgxpool.New(ctx, testUrl)
	if err != nil {
		logrus.Fatalf("error while trying to connect to test database: %v", err)
	}
	if err := testPool.Ping(ctx); err != nil {
		testPool.Close()
		logrus.Fatalf("error while trying to ping test database: %v", err)
	}

	// добавление миграций
	mgrt, err := migrate.New(migrationsPath, testUrl)
	if err != nil {
		logrus.Fatalf("error whyle trying to create migration object: %v", err)
	}
	if err := mgrt.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		mgrt.Close()
		logrus.Fatalf("error whyle trying to migrated tables: %v", err)
	}

	resultCode := m.Run()
	if err := mgrt.Down(); err != nil {
		mgrt.Close()
		logrus.Fatalf("error whyle trying to drop database: %v", err)
	}
	mgrt.Close()
	testPool.Close()
	os.Exit(resultCode)
}
