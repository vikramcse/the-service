package mid

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/vikramcse/the-service/internal/platform/web"
)

func Logger(log *log.Logger) web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(w, r)

			log.Printf("(%d) : %s %s -> %s (%s)",
				v.StatusCode, r.Method,
				r.URL.Path, r.RemoteAddr,
				time.Since(v.Start))

			return err
		}

		return h
	}

	return f
}
