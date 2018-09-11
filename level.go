package log

import (
	"strings"
)

const (
	// Trace is a log level used mainly for debugging purposes.
	// As this is the basic level, all other logs are emitted.
	Trace Level = iota

	// Info is a log level used to register relevant information about a process.
	// This level enables Warn, Error and Fatal levels too.
	Info

	// Warn is a log level that indicates some unexpected behavior that may not be a critical issue.
	// This level enables Error and Fatal levels too.
	Warn

	// Error is a log level that desired some attention. Usually indicated something very important.
	// This level enables the Fatal level too.
	Error

	// Fatal is a log level that indicates something very wrong, and also causes the application to halt.
	Fatal
)

// Level is used to categorize a logging message
type Level byte

func (l Level) String() string {
	switch l {
	case Trace:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	}
	return ""
}

// SetLevel configures the logging level by parsing the given string.
// The default level is Info.
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "trace", "debug":
		Default.Level = Trace
	case "warn":
		Default.Level = Warn
	case "error":
		Default.Level = Error
	case "fatal":
		Default.Level = Fatal
	default:
		Default.Level = Info
	}
}
