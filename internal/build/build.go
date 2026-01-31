package build

import (
	"chopper/internal/config"
	"chopper/internal/server"
)

func Run() error {
	serverConfig, err := config.ConfigsLoad()
	if err != nil {
		return err
	}

	server := server.NewServer(":8080", serverConfig.ReadTimeout, serverConfig.WriteTimeout, serverConfig.IdleTimeout, serverConfig.TimeToShutdown)
	if err := server.StartServer(); err != nil {
		return err
	}
	return nil
}
