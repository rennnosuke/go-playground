package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var schema = `
CREATE TABLE IF NOT EXISTS PERSON (
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255)
);`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// connect(test credential)
	userName := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")

	dsn := userName + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&autocommit=0&sql_mode=%27TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY%27"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// init
	db.MustExec(schema)

	// insert
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	// select
	var people []Person
	err = db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	if err != nil {
		panic(err)
	}
	fmt.Println(people)
}
