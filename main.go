package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	dbURL := "postgres://postgres:postgres@localhost:7777/inditilla?sslmode=disable"

	// Open a database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new PostgreSQL driver for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Force migration to a specific version
	versionToForce := 1 // replace with your desired version
	err = m.Force(versionToForce)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Migration forced to version %d\n", versionToForce)
}
