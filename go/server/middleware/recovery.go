package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/tehnerd/vatran/go/server/models"
)

// Recovery creates a middleware that recovers from panics.
// It logs the panic and stack trace, then returns a 500 Internal Server Error.
//
// Returns a middleware function that handles panics.
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic and stack trace
					log.Printf("panic: %v\n%s", err, debug.Stack())

					// Return internal server error
					models.WriteError(w, http.StatusInternalServerError,
						models.NewInternalError("internal server error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
