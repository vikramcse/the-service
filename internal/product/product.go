package product

import (
	"time"

	"github.com/google/uuid"
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

func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
			INSERT INTO products
			(product_id, name, cost, quantity, date_created, date_updated)
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &p, nil
}
