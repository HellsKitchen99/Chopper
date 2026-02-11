package config

import (
	"chopper/internal/domain"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func ConfigsLoad() (domain.ServerConfig, domain.JWtConfig, domain.DataBaseConfig, domain.RateLimiterConfig, error) {
	// загрузка конфига сервера
	var serverConfig domain.ServerConfig
	serverAddress := os.Getenv("SERVER_ADDRESS")
	serverReadtimeout := os.Getenv("SERVER_READTIMEOUT")
	serverWritetimeout := os.Getenv("SERVER_WRITETIMEOUT")
	serverIdletimeout := os.Getenv("SERVER_IDLETIMEOUT")
	serverTimeToShutdown := os.Getenv("SERVER_TIMETOSHUTDOWN")
	serverMode := os.Getenv("SERVER_MODE")
	if serverAddress == "" || serverReadtimeout == "" || serverWritetimeout == "" || serverIdletimeout == "" || serverTimeToShutdown == "" || serverMode == "" {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, fmt.Errorf("lack of environment variables")
	}
	readTimeout, err := time.ParseDuration(serverReadtimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	writeTimeout, err := time.ParseDuration(serverWritetimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	idleTimeout, err := time.ParseDuration(serverIdletimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	timeToShutdown, err := time.ParseDuration(serverTimeToShutdown)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	serverConfig.Address = serverAddress
	serverConfig.ReadTimeout = readTimeout
	serverConfig.WriteTimeout = writeTimeout
	serverConfig.IdleTimeout = idleTimeout
	serverConfig.TimeToShutdown = timeToShutdown
	var sm domain.ServerMode
	switch serverMode {
	case "release":
		sm = domain.ReleaseMode
	case "debug":
		sm = domain.DebugMode
	case "test":
		sm = domain.TestMode
	default:
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, fmt.Errorf("wrong server mode field")
	}
	serverConfig.ServerMode = sm

	// загрузка конфига базы данных
	var databaseConfig domain.DataBaseConfig
	databaseUser := os.Getenv("DB_USER")
	databasePassword := os.Getenv("DB_PASSWORD")
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	if databaseUser == "" || databasePassword == "" || databaseHost == "" || databasePort == "" || databaseName == "" {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, fmt.Errorf("lack of environment variables")
	}
	databaseConfig.User = databaseUser
	databaseConfig.Password = url.QueryEscape(databasePassword)
	databaseConfig.Host = databaseHost
	databaseConfig.Port = databasePort
	databaseConfig.Name = databaseName

	// загрузка конфига jwt сервиса
	var jwtConfig domain.JWtConfig
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpirationTime := os.Getenv("JWT_EXPIRATIONTIME")
	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtAudience := os.Getenv("JWT_AUDIENCE")
	if jwtSecret == "" || jwtExpirationTime == "" || jwtIssuer == "" || jwtAudience == "" {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, fmt.Errorf("lack of environment variables")
	}
	jwtConfig.Secret = []byte(jwtSecret)
	jwtValidatedExpirationTime, err := time.ParseDuration(jwtExpirationTime)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	jwtConfig.ExpirationTime = jwtValidatedExpirationTime
	jwtConfig.Issuer = jwtIssuer
	jwtConfig.Audience = jwtAudience

	// загрузка конфига рейт лимитера
	limiterRate := os.Getenv("LIMITER_RATE")
	limiterBurst := os.Getenv("LIMITER_BURST")
	if limiterRate == "" || limiterBurst == "" {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, fmt.Errorf("lack of environment variables")
	}
	parsedLimiterRate, err := time.ParseDuration(limiterRate)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	parsedLimiterBurst, err := strconv.Atoi(limiterBurst)
	if err != nil {
		return domain.ServerConfig{}, domain.JWtConfig{}, domain.DataBaseConfig{}, domain.RateLimiterConfig{}, err
	}
	var rateLimiterConfig domain.RateLimiterConfig
	rateLimiterConfig.Rate = parsedLimiterRate
	rateLimiterConfig.Burst = parsedLimiterBurst
	return serverConfig, jwtConfig, databaseConfig, rateLimiterConfig, nil
}
