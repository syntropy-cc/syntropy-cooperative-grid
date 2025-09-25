// Package middleware provides middleware components for the API
package middleware

import (
	"log"
	"os"
)

// Logger interface defines logging methods
type Logger interface {
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Fatal(message string, fields map[string]interface{})
}

// SimpleLogger implements a simple logger
type SimpleLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags),
		warn:  log.New(os.Stdout, "[WARN]  ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		fatal: log.New(os.Stderr, "[FATAL] ", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(message string, fields map[string]interface{}) {
	l.debug.Printf("%s %v", message, fields)
}

// Info logs an info message
func (l *SimpleLogger) Info(message string, fields map[string]interface{}) {
	l.info.Printf("%s %v", message, fields)
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(message string, fields map[string]interface{}) {
	l.warn.Printf("%s %v", message, fields)
}

// Error logs an error message
func (l *SimpleLogger) Error(message string, fields map[string]interface{}) {
	l.error.Printf("%s %v", message, fields)
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(message string, fields map[string]interface{}) {
	l.fatal.Printf("%s %v", message, fields)
	os.Exit(1)
}
