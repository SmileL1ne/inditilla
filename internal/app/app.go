package app

import (
	"context"
	"errors"
	"inditilla/config"
	"inditilla/internal/data"
	"inditilla/internal/handlers"
	"inditilla/internal/repository"
	"inditilla/internal/service"
	"inditilla/internal/service/user"
	"inditilla/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func Run(cfg *config.Config) {
	// Initialize new logger
	l, closeFile := logger.New(cfg.Log.Level)
	defer closeFile()

	// Open database connection
	db, err := openDB(cfg.Database.URL)
	if err != nil {
		l.Fatal(err.Error())
	}

	r := repository.New(db)
	deadline, err := strconv.Atoi(cfg.Auth.Deadline)
	if err != nil {
		l.Fatal(err.Error())
	}
	auth := user.NewAuthorizer([]byte(cfg.Auth.SigningKey), time.Duration(deadline)*time.Second)
	tokenModel := &data.TokenModel{Log: l}
	s := service.New(r, auth, tokenModel)
	fd := form.NewDecoder()

	logAdapter := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger().Level(zerolog.ErrorLevel)
	errLogger := log.New(logAdapter, "", 0)

	server := &http.Server{
		Addr:     "127.0.0.1:" + cfg.Http.Port,
		Handler:  handlers.NewRouter(l, s, fd),
		ErrorLog: errLogger,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		sig := <-sigCh
		l.Info("signal received: %s", sig.String())

		if err := server.Shutdown(context.Background()); err != nil {
			l.Fatal("server shutdown: %v", err)
		}
		if err := db.Close(context.Background()); err != nil {
			l.Fatal("db connection close: %v", err)
		}

		os.Exit(0)
	}()

	l.Info("starting the server: addr - %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Fatal("listen and serve: %v", err)
	}
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
