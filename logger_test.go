package logger_test

import (
	"fmt"
	"testing"

	"github.com/kyani-inc/logger"

	"github.com/stvp/go-udp-testing"
)

func TestWritingToPapertrail(t *testing.T) {
	port := 16661
	udp.SetAddr(fmt.Sprintf(":%d", port))

	log := logger.New(logger.Config{
		Appname: "myapp",
		Host:    "localhost",
		Port:    port,
	})

	udp.ShouldReceive(t, "foo", func() {
		log.Info("foo")
	})
}
