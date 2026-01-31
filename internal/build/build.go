package build

import (
	"chopper/internal/config"
	"chopper/internal/server"
	"time"
)

func Run() error {
	serverConfig, err := config.ConfigsLoad()
	if err != nil {
		return err
	}
	readTimeout, err := time.ParseDuration(serverConfig.ReadTimeout)
	if err != nil {
		return err
	}
	writeTimeout, err := time.ParseDuration(serverConfig.WriteTimeout)
	if err != nil {
		return err
	}
	idleTimeout, err := time.ParseDuration(serverConfig.IdleTimeout)
	if err != nil {
		return err
	}
	timeToShutdown, err := time.ParseDuration(serverConfig.TimeToShutdown)
	if err != nil {
		return err
	}
	server := server.NewServer(":8080", readTimeout, writeTimeout, idleTimeout, timeToShutdown)
	if err := server.StartServer(); err != nil {
		return err
	}
	return nil
}
