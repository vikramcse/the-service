package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * FROM products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

func Retrive(db *sqlx.DB, id string) (*Product, error) {
	var p Product

	const q = `SELECT * FROM products WHERE product_id = $1`
	if err := db.Get(&p, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting single product")
	}

	return &p, nil
}
