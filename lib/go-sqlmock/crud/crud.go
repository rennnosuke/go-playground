package crud

import (
	"context"
	"database/sql"
)

type Product struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}

func GetProducts(ctx context.Context, db *sql.DB) ([]Product, error) {
	if db == nil {
		return nil, nil
	}

	rows, err := db.QueryContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Product, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}