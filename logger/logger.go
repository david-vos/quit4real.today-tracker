package logger

import (
	"log"
	"os"
	"runtime"
	"strings"
)

// LogLevel represents the severity level of logging.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelFail
)

// CustomLogger wraps the standard logger and provides log level methods.
type CustomLogger struct {
	logger *log.Logger
	level  LogLevel
}

// Logger is the global logger instance
var Logger *CustomLogger

// NewCustomLogger creates a new CustomLogger with the specified log level.
func NewCustomLogger(level LogLevel) *CustomLogger {
	return &CustomLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		level:  level,
	}
}

// LogLevelFromEnv retrieves the log level from the environment.
func LogLevelFromEnv() LogLevel {
	levelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(levelStr) {
	case "debug":
		return LevelDebug
	case "prod":
		return LevelInfo // In prod mode, we log info and fail
	default:
		return LevelDebug // Default to debug level
	}
}

// InitLogger initializes the global logger
func InitLogger() {
	logLevel := LogLevelFromEnv()
	Logger = NewCustomLogger(logLevel)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info(msg)
}

// Fail logs a failure message
func Fail(msg string) {
	Logger.Fail(msg)
}

func (l *CustomLogger) Debug(msg string) {
	if l.level == LevelDebug {
		l.logger.Println("DEBUG:", msg)
	}
}

func (l *CustomLogger) Info(msg string) {
	if l.level == LevelDebug || l.level == LevelInfo {
		l.logger.Println("INFO:", msg)
	}
}

func (l *CustomLogger) Fail(msg string) {
	l.logger.Println("FAIL:", msg)
	l.logger.Println("Stack Trace:\n", captureStackTrace())
}

func captureStackTrace() string {
	stack := make([]byte, 1024)
	n := runtime.Stack(stack, true)
	return string(stack[:n])
}
