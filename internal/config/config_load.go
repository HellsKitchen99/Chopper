package config

import (
	"chopper/internal/domain"
	"fmt"
	"os"
	"time"
)

func ConfigsLoad() (domain.ServerConfig, error) {
	var serverConfig domain.ServerConfig
	serverAddress := os.Getenv("SERVER_ADDRESS")
	serverReadtimeout := os.Getenv("SERVER_READTIMEOUT")
	serverWritetimeout := os.Getenv("SERVER_WRITETIMEOUT")
	serverIdletimeout := os.Getenv("SERVER_IDLETIMEOUT")
	serverTimeToShutdown := os.Getenv("SERVER_TIMETOSHUTDOWN")
	if serverAddress == "" || serverReadtimeout == "" || serverWritetimeout == "" || serverIdletimeout == "" || serverTimeToShutdown == "" {
		return domain.ServerConfig{}, fmt.Errorf("lack of environment variables")
	}
	readTimeout, err := time.ParseDuration(serverReadtimeout)
	if err != nil {
		return domain.ServerConfig{}, err
	}
	writeTimeout, err := time.ParseDuration(serverWritetimeout)
	if err != nil {
		return domain.ServerConfig{}, err
	}
	idleTimeout, err := time.ParseDuration(serverIdletimeout)
	if err != nil {
		return domain.ServerConfig{}, err
	}
	timeToShutdown, err := time.ParseDuration(serverTimeToShutdown)
	if err != nil {
		return domain.ServerConfig{}, err
	}
	serverConfig.Address = serverAddress
	serverConfig.ReadTimeout = readTimeout
	serverConfig.WriteTimeout = writeTimeout
	serverConfig.IdleTimeout = idleTimeout
	serverConfig.TimeToShutdown = timeToShutdown
	return serverConfig, nil
}
