package crud

import (
	"context"
	"database/sql"
	"time"
)

// refer from squirrel
type StdSqlCtx interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

type Product struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func GetByID(ctx context.Context, db StdSqlCtx, id int) (*Product, error) {
	if db == nil {
		return nil, nil
	}
	row := db.QueryRowContext(ctx, "SELECT * FROM products WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var p Product
	if err := row.Scan(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func FindAll(ctx context.Context, db StdSqlCtx) ([]Product, error) {
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}

func Create(ctx context.Context, db StdSqlCtx, p *Product) (int64, error) {
	if db == nil {
		return 0, nil
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	result, err := db.ExecContext(ctx, "INSERT INTO products (name, price, created_at, updated_at) VALUES (?, ?, ?, ?)", p.Name, p.Price, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func Update(ctx context.Context, db StdSqlCtx, p *Product) error {
	if db == nil {
		return nil
	}
	p.UpdatedAt = time.Now()
	_, err := db.ExecContext(ctx, "UPDATE products SET name = ?, price = ?, updated_at = ? WHERE id = ?", p.Name, p.Price, p.UpdatedAt, p.ID)
	return err
}

func Upsert(ctx context.Context, db *sql.DB, p *Product) error {
	if db == nil {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			err = tx.Rollback()
		}
	}()

	var exists *Product
	if p.ID > 0 {
		exists, err = GetByID(ctx, tx, p.ID)
		if err != nil {
			return err
		}
	}
	if exists == nil {
		id, err := Create(ctx, tx, p)
		if err != nil {
			return err
		}
		p.ID = int(id)
		return nil
	}
	return Update(ctx, tx, p)
}

func Delete(ctx context.Context, db StdSqlCtx, id int) error {
	if db == nil {
		return nil
	}
	_, err := db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	return err
}
