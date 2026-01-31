package domain

import "time"

type ServerConfig struct {
	Address        string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	TimeToShutdown time.Duration
}
