package build

import (
	"chopper/internal/config"
	"chopper/internal/repositoty"
	"chopper/internal/server"
	"chopper/internal/usecase"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
)

var migrationsFilePath string = "file:///app/migrations"

func Run() error {
	fmt.Println("step1")
	// загрузка всех конфигов
	serverConfig, databaseConfig, err := config.ConfigsLoad()
	if err != nil {
		return err
	}

	fmt.Println(serverConfig, databaseConfig)
	fmt.Println("step2")
	// подключение к бд
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", databaseConfig.User, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.Name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	defer pool.Close()

	fmt.Println("step3")
	// миграции
	m, err := migrate.New(migrationsFilePath, connStr)
	fmt.Println("MIGRATE RESULT:", err)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	fmt.Println("step4")
	// создание слоев
	userRepo := repositoty.NewUserRepositoryRealization(pool)
	userService := usecase.NewUserService(userRepo)

	fmt.Println("step5")
	// запуск сервера
	server := server.NewServer(":8080", serverConfig.ReadTimeout, serverConfig.WriteTimeout, serverConfig.IdleTimeout, serverConfig.TimeToShutdown, userService)
	if err := server.StartServer(); err != nil {
		return err
	}
	return nil
}
