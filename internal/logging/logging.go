package logging

import (
	"fmt"
	"io"
	"time"
)

// LogLevel defines the level of a log message.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// Logger defines the interface for a logging backend.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// DefaultLogger is a simple logger that logs messages to standard output.
type DefaultLogger struct {
	Out io.Writer
}

// NewDefaultLogger creates a new DefaultLogger that logs messages to the given writer.
func NewDefaultLogger(out io.Writer) *DefaultLogger {
	return &DefaultLogger{Out: out}
}

// Debugf logs a debug message.
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.log(LogLevelDebug, format, v...)
}

// Infof logs an informational message.
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.log(LogLevelInfo, format, v...)
}

// Warningf logs a warning message.
func (l *DefaultLogger) Warningf(format string, v ...interface{}) {
	l.log(LogLevelWarning, format, v...)
}

// Errorf logs an error message.
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.log(LogLevelError, format, v...)
}

// log logs a message with the given level.
func (l *DefaultLogger) log(level LogLevel, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.Out, "[%s] [%s] %s\n", timestamp, levelToString(level), message)
}

// levelToString returns the string representation of a LogLevel.
func levelToString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
