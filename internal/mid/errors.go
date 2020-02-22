package mid

import (
	"log"
	"net/http"

	"github.com/vikramcse/the-service/internal/platform/web"
)

// Errors handles erros coming out of the call chanin. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged
func Errors(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if err := before(w, r); err != nil {
				log.Printf("ERROR: %+v", err)

				if err := web.RespondError(w, err); err != nil {
					return err
				}
			}

			return nil
		}

		return h
	}

	return f
}
