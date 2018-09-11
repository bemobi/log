package log

import (
	"io"
	"os"
	"time"
)

func init() {
	Default = &Emitter{Level: Info, Output: os.Stderr}
}

// Default is the default logger engine
var Default *Emitter

// SetWithTime configures a logger with timestamp
func SetWithTime() {
	Default = &Emitter{Level: Info, Output: os.Stderr, TimeFormat: time.RFC3339}
}

// SetHook configures a logger with an emitter hook
func SetHook(hook Hook) {
	Default = &Emitter{
		Level:  Info,
		Output: os.Stderr,
		Hook:   hook,
	}
}

// T logs a formatted message when the Level is set to Trace
func T(tag string, message string, v ...interface{}) {
	Default.Emit(tag, Trace, message, v...)
}

// I logs a formatted message when the Level is set to Info or lower
func I(tag string, message string, v ...interface{}) {
	Default.Emit(tag, Info, message, v...)
}

// W logs a formatted message when the Level is set to Warn or lower
func W(tag string, message string, v ...interface{}) {
	Default.Emit(tag, Warn, message, v...)
}

// E logs a formatted message when the Level is set to Error or lower
func E(tag string, message string, v ...interface{}) {
	Default.Emit(tag, Error, message, v...)
}

// F logs a formatted message when the Level is set to Fatal or lower
func F(tag string, message string, v ...interface{}) {
	Default.Emit(tag, Fatal, message, v...)
	_exit(1)
}

// SetTestMode toggles the testing mode, which is disabled by default.
//
// When testing mode is on, all the logging functions emit values to the first sink (io.Writer)
// and the program does not halt when logging with F.
//
// If you need to test the exit operation when logging with F, check the _exitCode global variable.
func SetTestMode(active bool, sink ...io.Writer) {
	if active {
		// change the exit function to generate panics
		_exit = func(code int) {
			_exitCode = code
		}
		// backup the current logger
		_logger = Default
		// create a new logger
		Default = &Emitter{
			Output: sink[0],
			Hook:   _logger.Hook,
		}
	} else {
		_exit = os.Exit
		Default = _logger
	}
}

// LastExitCode returns the last exit code produced by a Fatal logging
//
// For testing purposes only
func LastExitCode() int {
	return _exitCode
}

var (
	// For tests
	_exit     = os.Exit
	_logger   = Default
	_exitCode = 0
)
