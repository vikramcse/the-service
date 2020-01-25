package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/vikramcse/the-service/internal/platform/web"
)

func API(db *sqlx.DB, log *log.Logger) http.Handler {
	app := web.NewApp(log)
	p := Products{DB: db, Log: log}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrive)
	app.Handle(http.MethodPost, "/v1/products", p.Create)

	return app
}
