package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

// Initialize database migration up before start of the server
func init() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Str("state", "error loading '.env' file").Msg("migrate")
	}

	dbURL, ok := os.LookupEnv("DB_URL") // Url in .env is empty for now
	if !ok || len(dbURL) == 0 {
		logger.Fatal().Str("state", "environment variable not declared: DB_URL").Msg("migrate")
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

		logger.Warn().Str("state", fmt.Sprintf("postgres is trying to connect, attempts left: %d", attempts)).Msg("migrate")
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Fatal().Err(err).Str("state", "postgres connection").Msg("migrate")
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal().Err(err).Str("state", "up").Msg("migrate")
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Warn().Str("state", "no change").Msg("migrate")
		return
	}

	logger.Info().Str("migrate", "success")
}
