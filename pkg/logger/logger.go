package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

// This ensures that Logger struct implements ILogger interface
var _ ILogger = (*Logger)(nil)

// New returns new logger with specified level and file closing function
// that should be defered when this function is called
func New(lvl string) (*Logger, func() error) {
	var l zerolog.Level

	switch strings.ToLower(lvl) {
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	// Get log file full path to open/create log file to save logs
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting working directory: %v", err)
	}
	logFilePath := "logs.txt"
	fullLogFilePath := filepath.Join(curDir, logFilePath)

	// Open/Create file to save logs (where also user activity is tracked)
	fileWriter, err := os.OpenFile(fullLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// Create new logger with 2 outputs - console and log file
	multi := zerolog.MultiLevelWriter(zerolog.NewConsoleWriter(), fileWriter)
	logger := zerolog.New(multi).With().Timestamp().Caller().Logger()

	return &Logger{
		logger: &logger,
	}, fileWriter.Close
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.msg("error", message, args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.logger.Error().Msgf(msg.Error(), args...)
	case string:
		l.logger.Error().Msgf(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unkown type %v", level, message, msg), args...)
	}
}
