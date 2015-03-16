package logger

import (
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/polds/logrus/hooks/papertrail"
)

var __l *logrus.Logger

// Logger returns an instance of
// a logger or creates a new one.
func Logger() *logrus.Logger {
	if __l == nil {
		__l = NewLogger()
	}

	return __l
}

// NewLogger creates a new instances of a
// logrus logger and returns the instance.
func NewLogger() *logrus.Logger {
	app := os.Getenv("APPNAME")
	if app == "" {
		app, _ = os.Hostname()
		os.Setenv("APPNAME", app)
	}

	host := os.Getenv("PAPERTRAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))

	log := logrus.New()
	hook, err := logrus_papertrail.NewPapertrailHook(host, port, app)
	hook.UseHostname()

	// Register the PaperTrail hook
	if err == nil {
		log.Hooks.Add(hook)
	}

	return log
}
