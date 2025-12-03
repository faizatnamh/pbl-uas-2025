package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// TEMPORARY DEBUG:
	fmt.Println("HOST:", host)
	fmt.Println("PORT:", port)
	fmt.Println("USER:", user)
	fmt.Println("PASS:", password)
	fmt.Println("DB:", dbname)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL")
	return db
}
