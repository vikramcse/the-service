package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vikramcse/the-service/internal/platform/web"
	"github.com/vikramcse/the-service/internal/product"
)

type Products struct {
	DB  *sqlx.DB
	Log *log.Logger
}

func (p *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, list, http.StatusOK)
}

func (p *Products) Retrive(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrive(r.Context(), p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}
	return web.Respond(w, prod, http.StatusOK)
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct

	if err := json.NewDecoder(r.Body).Decode(&np); err != nil {
		return errors.Wrapf(err, "decoding new product")
	}

	prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, &prod, http.StatusCreated)
}
