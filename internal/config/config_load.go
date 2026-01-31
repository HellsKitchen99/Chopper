package config

import (
	"chopper/internal/domain"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func ConfigsLoad() (domain.ServerConfig, error) {
	if err := godotenv.Load(); err != nil {
		return domain.ServerConfig{}, err
	}
	var serverConfig domain.ServerConfig
	if err := env.Parse(&serverConfig); err != nil {
		return domain.ServerConfig{}, err
	}

	return serverConfig, nil
}
