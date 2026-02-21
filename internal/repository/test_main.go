package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var testPool *pgxpool.Pool

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
		testPool.Close()
		logrus.Fatalf("error while trying to connect to test database: %v", err)
	}
	if err := testPool.Ping(ctx); err != nil {
		testPool.Close()
		logrus.Fatalf("error while trying to ping test database: %v", err)
	}
	resultCode := m.Run()
	testPool.Close()
	os.Exit(resultCode)
}
