package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key
type ctxKey int

// keyValues is how request values are stored/retrived
const KeyValues ctxKey = 1

// Values carries information about reach request
type Values struct {
	StatusCode int
	Start      time.Time
}

type Handler func(http.ResponseWriter, *http.Request) error

// App is the entrypoint to this application and what controls the context of
// each request.
type App struct {
	log *log.Logger
	mux *chi.Mux
	mw  []Middleware
}

// NewApp constructs an App to handle a set of routes. Any Middleware provided
// will be ran for every request
func NewApp(log *log.Logger, mw ...Middleware) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
		mw:  mw,
	}
}

// Handle associates a handler function with an HTTP method and URL Pattern

// It converts our custom handle type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way
func (a *App) Handle(method, url string, h Handler) {
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)
		r = r.WithContext(ctx)

		if err := h(w, r); err != nil {
			a.log.Printf("Unhandled error: %+v", err)
		}
	}
	a.mux.MethodFunc(method, url, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
