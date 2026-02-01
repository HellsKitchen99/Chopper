package config

import (
	"chopper/internal/domain"
	"fmt"
	"net/url"
	"os"
	"time"
)

func ConfigsLoad() (domain.ServerConfig, domain.DataBaseConfig, error) {
	var serverConfig domain.ServerConfig
	serverAddress := os.Getenv("SERVER_ADDRESS")
	serverReadtimeout := os.Getenv("SERVER_READTIMEOUT")
	serverWritetimeout := os.Getenv("SERVER_WRITETIMEOUT")
	serverIdletimeout := os.Getenv("SERVER_IDLETIMEOUT")
	serverTimeToShutdown := os.Getenv("SERVER_TIMETOSHUTDOWN")
	if serverAddress == "" || serverReadtimeout == "" || serverWritetimeout == "" || serverIdletimeout == "" || serverTimeToShutdown == "" {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, fmt.Errorf("lack of environment variables")
	}
	readTimeout, err := time.ParseDuration(serverReadtimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, err
	}
	writeTimeout, err := time.ParseDuration(serverWritetimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, err
	}
	idleTimeout, err := time.ParseDuration(serverIdletimeout)
	if err != nil {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, err
	}
	timeToShutdown, err := time.ParseDuration(serverTimeToShutdown)
	if err != nil {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, err
	}
	serverConfig.Address = serverAddress
	serverConfig.ReadTimeout = readTimeout
	serverConfig.WriteTimeout = writeTimeout
	serverConfig.IdleTimeout = idleTimeout
	serverConfig.TimeToShutdown = timeToShutdown
	var databaseConfig domain.DataBaseConfig
	databaseUser := os.Getenv("DB_USER")
	databasePassword := os.Getenv("DB_PASSWORD")
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	if databaseUser == "" || databasePassword == "" || databaseHost == "" || databasePort == "" || databaseName == "" {
		return domain.ServerConfig{}, domain.DataBaseConfig{}, fmt.Errorf("lack of environment variables")
	}
	databaseConfig.User = databaseUser
	databaseConfig.Password = url.QueryEscape(databasePassword)
	databaseConfig.Host = databaseHost
	databaseConfig.Port = databasePort
	databaseConfig.Name = databaseName
	return serverConfig, databaseConfig, nil
}
