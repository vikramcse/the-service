package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/vikramcse/the-service/internal/mid"
	"github.com/vikramcse/the-service/internal/platform/web"
)

func API(db *sqlx.DB, log *log.Logger) http.Handler {
	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		p := Products{DB: db, Log: log}

		app.Handle(http.MethodGet, "/v1/products", p.List)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrive)
		app.Handle(http.MethodPost, "/v1/products", p.Create)

		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)
	}

	return app
}
