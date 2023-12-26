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

func Create(ctx context.Context, db *sql.DB, p Product) (int64, error) {
	if db == nil {
		return 0, nil
	}
	result, err := db.ExecContext(ctx, "INSERT INTO products (name, price) VALUES (?, ?)", p.Name, p.Price)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func Update(ctx context.Context, db *sql.DB, p Product) error {
	if db == nil {
		return nil
	}
	_, err := db.ExecContext(ctx, "UPDATE products SET name = ?, price = ? WHERE id = ?", p.Name, p.Price, p.ID)
	return err
}
