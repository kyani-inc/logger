package logger

import (
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/polds/logrus/hooks/papertrail"
)

var __l *logrus.Logger

type Config struct {
	Appname string
	Host    string
	Port    int

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
func New(config Config) Config {
	config.Logger = logrus.New()
	hook, err := logrus_papertrail.NewPapertrailHook(config.Host, config.Port, config.Appname)
	hook.UseHostname()

	// Register the PaperTrail hook
	if err == nil {
		config.Logger.Hooks.Add(hook)
	}

	return config
}

// DefaultConfig makes certain assumptions about your environment variables
// and can be used to create a new basic instance.
func DefaultConfig() Config {
	app := os.Getenv("APPNAME")
	if app == "" {
		app, _ = os.Hostname()
		os.Setenv("APPNAME", app)
	}

	port, _ := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))
	return Config{
		Appname: os.Getenv("APPNAME"),
		Host:    os.Getenv("PAPERTRAIL_HOST"),
		Port:    port,
	}
}

// @TODO polds deprecate this function
func NewLogger() *logrus.Logger {
	return New(DefaultConfig()).Logger
}
