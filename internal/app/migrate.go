package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	logger := zerolog.New(zerolog.NewConsoleWriter())

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("migrate - error loading .env")
	}

	dbURL, ok := os.LookupEnv("DB_URL") // Url in .env is empty for now
	if !ok || len(dbURL) == 0 {
		logger.Fatal().Str("migrate", "environment variable not declared: DB_URL")
	}

	sslMode, ok := os.LookupEnv("DB_SSL_MODE")
	if !ok || len(sslMode) == 0 {
		sslMode = "disable"
	}

	dbURL += "?sslmode=" + sslMode

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbURL)
		if err == nil {
			break
		}

		logger.Warn().Str("migrate", fmt.Sprintf("postgres is trying to connect, attempts left: %d", attempts))
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("migrate - postgres connection error")
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal().Err(err).Msg("migrate - up error")
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Warn().Str("migrate", "no change")
		return
	}

	logger.Info().Str("migrate", "success")
}
