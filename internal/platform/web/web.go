package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is the entrypoint to this application and what controls the context of
// each request.
type App struct {
	log *log.Logger
	mux *chi.Mux
}

// NewApp constructs an App to handle a set of routes
func NewApp(log *log.Logger) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
	}
}

// Handle associates a handler function with an HTTP method and URL Pattern
func (a *App) Handle(method, url string, h http.HandlerFunc) {
	a.mux.MethodFunc(method, url, h)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
