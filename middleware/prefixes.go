package middleware

import (
	"fmt"

	"github.com/go-martini/martini"
	"github.com/labstack/echo"
)

func makePrefixes(c interface{}, prefixes ...interface{}) (prefix string) {
	if len(prefixes) == 0 {
		return
	}

	for _, pfx := range prefixes {
		switch pfx := pfx.(type) {
		case string:
			prefix += fmt.Sprintf("[%s]", pfx)
		case func(c *echo.Context) string:
			prefix += fmt.Sprintf("[%s]", pfx(c.(*echo.Context)))
		case func(c martini.Context) string:
			prefix += fmt.Sprintf("[%s]", pfx(c.(martini.Context)))
		default:
			prefix += fmt.Sprintf("[%v]", pfx)
		}
	}

	prefix += " "
	return
}
