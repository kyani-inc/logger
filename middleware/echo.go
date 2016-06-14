package middleware

import (
	"net/http"
	"time"

	"github.com/kyani-inc/go-utils/ip"
	"github.com/kyani-inc/logger"

	"github.com/kyani-inc/echo"
)

func Echo(log logger.Client, prefixes ...interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			// Get a single IP Address of the connecting party
			addr := ip.Client(c.Request())

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)
			// Create the prefix for the logger
			prefix := makePrefixes(c, prefixes...)
			log.Infof("%s%v %s %s %v %s \"%s\"", prefix, c.Response().Status(), http.StatusText(c.Response().Status()), c.Request().Method, latency, c.Request().URL.Path, addr)

			return nil
		}
	}
}
