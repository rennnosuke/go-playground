package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var createPerson = `
CREATE TABLE IF NOT EXISTS PERSON (
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);
`

var createPlace = `
CREATE TABLE IF NOT EXISTS PLACE (
	country VARCHAR(255) NOT NULL,
    city VARCHAR(255),
    telcode INT NOT NULL
);
`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
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
	db.MustExec(createPerson)
	db.MustExec(createPlace)

	ctx := context.Background()

	// insert
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES (?, ?)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES (?, ?)", "Singapore", "65")
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	// select
	var people []Person
	err = db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[db.Select]: %s\n", people)

	// select with context
	var people2 []Person
	err = db.SelectContext(ctx, &people2, "SELECT * FROM person ORDER BY first_name ASC")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[db.SelectContext]: %s\n", people2)

	// get
	var person Person
	err = db.Get(&person, "SELECT * FROM person WHERE first_name=?", "Jason")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[db.Get]: %s\n", person)

	// get with context
	var person2 Person
	err = db.GetContext(ctx, &person2, "SELECT * FROM person WHERE first_name=?", "John")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[db.GetContext]: %s\n", person2)

	place := Place{}
	rows, err := db.Queryx("SELECT * FROM place")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", place)
	}
}
