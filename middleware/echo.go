package middleware

import (
	"net/http"
	"time"

	"github.com/kyani-inc/go-utils/ip"
	"gopkg.in/kyani-inc/logger.v2"

	"github.com/labstack/echo"
)

func Echo(log logger.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			// Get a single IP Address of the connecting party
			addr := ip.Client(c.Request())

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)
			log.Infof("%v %s %s %v %s \"%s\"", c.Response().Status(), http.StatusText(c.Response().Status()), c.Request().Method, latency, c.Request().URL.Path, addr)

			return nil
		}
	}
}
