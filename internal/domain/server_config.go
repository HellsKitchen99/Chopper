package domain

type ServerConfig struct {
	Address        string `env:"SERVER_ADDRESS,required"`
	ReadTimeout    string `env:"SERVER_READTIMEOUT,required"`
	WriteTimeout   string `env:"SERVER_WRITETIMEOUT,required"`
	IdleTimeout    string `env:"SERVER_IDLETIMEOUT,required"`
	TimeToShutdown string `env:"SERVER_TIMETOSHUTDOWN,required"`
}
