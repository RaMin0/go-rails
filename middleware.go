package rails

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

var middlewares = middleware{
	// wrapPanic
	func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, fmt.Sprintf("[Panic] %+v", err), http.StatusInternalServerError)
				}
			}()

			fn(w, r)
		}
	},

	// withLog
	func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c := httptest.NewRecorder()
			fn(c, r)
			log.Printf("[%d] %s %s as %s\n",
				c.Code, r.Method, r.URL.Path,
				strings.SplitN(c.Header().Get("Content-Type"), ";", 2)[0])

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}
			w.WriteHeader(c.Code)
			c.Body.WriteTo(w)
		}
	},
}

type middleware []func(http.HandlerFunc) http.HandlerFunc

func (mws *middleware) wrap(fn http.HandlerFunc) http.HandlerFunc {
	for _, mw := range *mws {
		fn = mw(fn)
	}
	return fn
}
