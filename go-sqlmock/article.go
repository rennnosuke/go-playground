package sqlmock

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Article struct {
	ID      int
	Title   string
	Content string
}

func GetByID(id int, db *sql.DB) (*Article, error) {
	row := db.QueryRow("SELECT * FROM ARTICLES WHERE ID = ? AND IS_DELETED = 0", id)

	e := Article{}
	if err := row.Scan(&e.ID, &e.Title, &e.Content); err != nil {
		return nil, fmt.Errorf("failed to scan row: %s", err)
	}

	return &Article{ID: e.ID, Title: e.Title, Content: e.Content}, nil
}

func Create(id int, title, content string, db *sql.DB) error {
	tx, err := db.Begin()
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO ARTICLES (ID, TITLE, CONTENT) VALUES (?, ?, ?)", id, title, content)
	if err != nil {
		return fmt.Errorf("failed to insert article: %s", err)
	}

	return nil
}

func main() {
	db, err := sql.Open("mysql", "root@/blogsystem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db")
		os.Exit(1)
	}

	a, err := GetByID(1, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get article: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%v", a)
}
