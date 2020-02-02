package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("Product not found")
var ErrInvalidID = errors.New("ID is not in it's proper form")

func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `
			SELECT 
				p.*,
				COALESCE(SUM(s.quantity), 0) as sold,
				COALESCE(SUM(s.paid), 0) as revenue
			FROM products as p
			LEFT JOIN sales as s ON(p.product_id=s.product_id)
			GROUP BY p.product_id`

	if err := db.SelectContext(ctx, &products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

func Retrive(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product
	const q = `
			SELECT 
				p.*,
				COALESCE(SUM(s.quantity), 0) as sold,
				COALESCE(SUM(s.paid), 0) as revenue
			FROM products as p
			LEFT JOIN sales as s ON(p.product_id=s.product_id)
			WHERE p.product_id = $1
			GROUP BY p.product_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
	}

	return &p, nil
}

func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
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

	_, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &p, nil
}
