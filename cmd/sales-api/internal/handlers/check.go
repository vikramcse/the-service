package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/vikramcse/the-service/internal/platform/database"
	"github.com/vikramcse/the-service/internal/platform/web"
)

type Check struct {
	db *sqlx.DB
}

func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.db); err != nil {
		health.Status = "db not ready"
		return web.Respond(w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(w, health, http.StatusOK)
}
