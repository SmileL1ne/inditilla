package app

import (
	"context"
	"inditilla/config"
	"inditilla/pkg/logger"

	"github.com/jackc/pgx/v5"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	_, err := openDB(cfg.Database.URL) // Get database connection into 'pkg'
	if err != nil {
		l.Fatal(err.Error())
	}

	// Create repository
	// Create service

	// Initialize server

	// Graceful shutdown here

	// Start the server
}

func openDB(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, err
}
