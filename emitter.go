package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

// New returns a new logging emmitter
//
// The standard output is os.Stderr.
//
// The most used time format is time.RFC3339.
// However, if your logs are shipped via syslog, you can omit the time format.
func New(output io.Writer, timeFormat string) *Emitter {
	return &Emitter{Output: output, TimeFormat: timeFormat}
}

// Emitter is the base logging type
type Emitter struct {
	Level      Level
	Output     io.Writer
	TimeFormat string
	Hook       Hook

	// parsed fields
	context []byte
}

// T logs a formatted message when the Level is set to Trace
func (e *Emitter) T(tag string, message string, fields ...interface{}) {
	e.Emit(tag, Trace, message, fields...)
}

// I logs a formatted message when the Level is set to Info or lower
func (e *Emitter) I(tag string, message string, fields ...interface{}) {
	e.Emit(tag, Info, message, fields...)
}

// W logs a formatted message when the Level is set to Warn or lower
func (e *Emitter) W(tag string, message string, fields ...interface{}) {
	e.Emit(tag, Warn, message, fields...)
}

// E logs a formatted message when the Level is set to Error or lower
func (e *Emitter) E(tag string, message string, fields ...interface{}) {
	e.Emit(tag, Error, message, fields...)
}

// F logs a formatted message when the Level is set to Fatal or lower
func (e *Emitter) F(tag string, message string, fields ...interface{}) {
	e.Emit(tag, Fatal, message, fields...)
	_exit(1)
}

// Emit formats and writes a logging message to the emitters' output
func (e *Emitter) Emit(tag string, level Level, message string, fields ...interface{}) {
	if level < e.Level {
		return
	}

	buf := pool.Get().(*bytes.Buffer)

	// start document
	buf.WriteByte('{')

	// time
	if e.TimeFormat != "" {
		buf.WriteString(`"time":"`)
		buf.WriteString(time.Now().Format(e.TimeFormat))
		buf.WriteString(`",`)
	}

	// tag
	buf.WriteString(`"tag":"`)
	buf.WriteString(tag)
	buf.WriteString(`",`)

	// level
	buf.WriteString(`"level":`)
	buf.WriteByte('"')
	buf.WriteString(level.String())
	buf.WriteByte('"')

	// message
	buf.WriteString(`,"msg":`)
	writeJSONString(buf, message)

	// fields
	if e.context != nil {
		buf.Write(e.context)
	}
	writeFields(buf, fields...)

	// end document
	buf.WriteByte('}')
	buf.WriteByte('\n')

	// call hook
	if e.Hook != nil {
		b := buf.Bytes()
		e.Hook(level, b)
	}

	buf.WriteTo(e.Output)
	pool.Put(buf)
}

// Hook defines an emitter hook
type Hook func(Level, []byte)

func writeFields(buf *bytes.Buffer, fields ...interface{}) {
	for field := 0; field < len(fields); field += 2 {
		// WriteString is slower in this codepath
		buf.WriteByte(',')

		// Key
		buf.WriteByte('"')
		switch k := fields[field].(type) {
		case string:
			buf.WriteString(k)
		case fmt.Stringer:
			buf.WriteString(k.String())
		default:
			fmt.Fprintf(buf, `%v`, k)
		}

		buf.WriteString(`":`)

		// Value
		switch val := fields[field+1].(type) {
		case string:
			writeJSONString(buf, val)
		case []byte:
			writeJSONString(buf, string(val))
		case fmt.Stringer:
			writeJSONString(buf, val.String())
		case byte, int, int8, int16, int32, int64, float32, float64, bool, uint, uint16, uint32, uint64:
			fmt.Fprintf(buf, `%v`, val)
		case error:
			writeJSONString(buf, val.Error())
		default:
			writeJSONString(buf, fmt.Sprintf(`%v`, val))
		}
	}
}

func writeJSONString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch b {
		case '"', '\\':
			buf.WriteByte('\\')
			buf.WriteByte(b)
		case '\b':
			buf.WriteByte('\\')
			buf.WriteByte('b')
		case '\f':
			buf.WriteByte('\\')
			buf.WriteByte('f')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\t':
			buf.WriteByte('\\')
			buf.WriteByte('t')
		default:
			buf.WriteByte(b)
		}
	}
	buf.WriteByte('"')
}

var pool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}
