package logger

import (
	"os"
)

type Config struct {
	Appname string
	Host    string
	Port    int
}

type Client struct {
	*Logger
}

// New creates a new instance of a
// logger based on provided config data
func New(config Config) Client {
	var client Client

	client.Logger = NewLogger(config)

	return client
}

// DefaultConfig makes certain assumptions about your environment variables
// and can be used to create a new basic instance.
func DefaultConfig() Config {
	app := os.Getenv("APPNAME")
	if app == "" {
		app, _ = os.Hostname()
	}
	return Config{
		Appname: app,
		Host:    os.Getenv("SUMO_ENDPOINT"),
	}
}

func NewLoggerSumo() *Logger {
	return NewSumo(DefaultConfig()).Logger
}
