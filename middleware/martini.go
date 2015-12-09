package middleware

import (
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/kyani-inc/go-utils/ip"
	"github.com/kyani-inc/logger"
)

// Martini is a Martini Middleware that emulates their default logger in
// the sense that it logs every request and sends it to Papertrail.
// The logger gets the following info:
// 	- Time taken for request
//	- IP Address of connecting party
//	- HTTP Status Code
// 	- HTTP Status Text
// 	- Request Method
// 	- Request Path
//
// The result appears in Papertrail as:
// [info] 200 OK HEAD 1.109259ms /my/endpoint/ "8.8.8.8"
func Martini(prefixes ...interface{}) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log logger.Client) {
		start := time.Now()

		// Get a single IP Address of the connecting party
		addr := ip.Client(req)

		rw := res.(martini.ResponseWriter)
		c.Next()

		prefix := makePrefixes(c, prefixes...)
		log.Infof("%s%v %s %s %v %s \"%s\"", prefix, rw.Status(), http.StatusText(rw.Status()), req.Method, time.Since(start), req.URL.Path, addr)
	}
}
