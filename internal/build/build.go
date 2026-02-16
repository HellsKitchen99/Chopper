package build

import (
	"chopper/internal/config"
	"chopper/internal/middleware"
	"chopper/internal/repository"
	"chopper/internal/security"
	"chopper/internal/server"
	"chopper/internal/usecase"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/time/rate"

	"github.com/jackc/pgx/v5/pgxpool"
)

var migrationsFilePath string = "file:///app/migrations"

func Run() error {
	fmt.Println("step1")
	// загрузка всех конфигов
	serverConfig, jwtConfig, databaseConfig, rateLimiterConfig, err := config.ConfigsLoad()
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
	userRepo := repository.NewUserRepositoryRealization(pool)
	jwtService := security.NewJwt(jwtConfig.Secret, jwtConfig.ExpirationTime, jwtConfig.Issuer, jwtConfig.Audience)
	// сделать конфиг с кост
	cost := 10
	passwordHasher := security.NewPasswordHasher(cost)
	uuidGenerator := security.NewUUIDGenerator()
	userService := usecase.NewUserService(userRepo, jwtService, passwordHasher, uuidGenerator)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	dailyNotesRepo := repository.NewDailyNotesRepositoryRealization(pool)
	dailyNotesService := usecase.NewDailyNotesService(dailyNotesRepo)
	alertRepository := repository.NewAlertRepositoryRealization(pool)
	alertService := usecase.NewAlertServcie(alertRepository)
	rateLimiter := middleware.NewRateLimiter(rate.Every(rateLimiterConfig.Rate), rateLimiterConfig.Burst)

	fmt.Println("step5")
	// запуск сервера
	server := server.NewServer(serverConfig.Address, serverConfig.ReadTimeout, serverConfig.WriteTimeout, serverConfig.IdleTimeout, serverConfig.TimeToShutdown, serverConfig.ServerMode, userService, dailyNotesService, alertService, authMiddleware, rateLimiter)
	if err := server.StartServer(); err != nil {
		return err
	}
	return nil
}
