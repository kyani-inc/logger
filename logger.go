package logger

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/kyani-inc/logrus-papertrail-hook.v3"
)

var __l *logrus.Logger

type Config struct {
	Appname string
	Host    string
	Port    int
}

type Client struct {
	*logrus.Logger
}

// Logger returns an instance of
// a logger or creates a new one.
func Logger() *logrus.Logger {
	if __l == nil {
		__l = NewLogger()
	}

	return __l
}

// New creates a new instance of a
// logger based on provided config data
func New(config Config) Client {
	var client Client

	host, _ := os.Hostname()

	client.Logger = logrus.New()
	client.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	hook, err := logrus_papertrail.NewPapertrailHook(&logrus_papertrail.Hook{
		Host:     config.Host,
		Port:     config.Port,
		Appname:  config.Appname,
		Hostname: host,
	})

	// Register the PaperTrail hook
	if err == nil {
		client.Logger.Hooks.Add(hook)
	}

	return client
}

// DefaultConfig makes certain assumptions about your environment variables
// and can be used to create a new basic instance.
func DefaultConfig() Config {
	app := os.Getenv("APPNAME")
	if app == "" {
		app, _ = os.Hostname()
	}
	if os.Getenv("PAPERTRAIL_PORT") != "" {
		port, _ := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))
		return Config{
			Appname: app,
			Host:    os.Getenv("PAPERTRAIL_HOST"),
			Port:    port,
		}
	}

	return Config{
		Appname: app,
		Host:    os.Getenv("SUMO_ENDPOINT"),
	}
}

// @TODO polds deprecate this function
func NewLogger() *logrus.Logger {
	return New(DefaultConfig()).Logger
}

func NewLoggerSumo() *logrus.Logger {
	return NewSumo(DefaultConfig()).Logger
}
