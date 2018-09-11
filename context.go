package log

import (
	"bytes"
)

// C creates a logging context
func C(tag string, fields ...interface{}) *Context {
	var buf bytes.Buffer
	writeFields(&buf, fields...)
	emitter := &Emitter{
		Level:      Default.Level,
		Output:     Default.Output,
		TimeFormat: Default.TimeFormat,
		Hook:       Default.Hook,
		context:    buf.Bytes(),
	}
	return &Context{Emitter: emitter, Tag: tag}
}

// Context is a logging emitter wrapper with parsed context fields
//
// The common use case is
//
//		logger := log.C("TAG", "one", 1, "two", 2)
// 		logger.I("start")
// 		logger.I("stop")
//
// The output is something like below
//
//		{"tag":"TAG","msg":"start","one":1,"two":2}
//		{"tag":"TAG","msg":"stop","one":1,"two":2}
type Context struct {
	Emitter *Emitter
	Tag     string
}

// C returns a new context based on the current context
func (c *Context) C(fields ...interface{}) *Context {
	var buf bytes.Buffer
	buf.Write(c.Emitter.context)
	writeFields(&buf, fields...)
	emitter := &Emitter{
		Level:      c.Emitter.Level,
		Output:     c.Emitter.Output,
		TimeFormat: c.Emitter.TimeFormat,
		Hook:       c.Emitter.Hook,
		context:    buf.Bytes(),
	}
	return &Context{Emitter: emitter, Tag: c.Tag}
}

// T logs a message when the Level is set to Trace
func (c *Context) T(message string, fields ...interface{}) {
	c.Emitter.Emit(c.Tag, Trace, message, fields...)
}

// I logs a message when the Level is set to Info or lower
func (c *Context) I(message string, fields ...interface{}) {
	c.Emitter.Emit(c.Tag, Info, message, fields...)
}

// W logs a message when the Level is set to Warn or lower
func (c *Context) W(message string, fields ...interface{}) {
	c.Emitter.Emit(c.Tag, Warn, message, fields...)
}

// E logs a message when the Level is set to Error or lower
func (c *Context) E(message string, fields ...interface{}) {
	c.Emitter.Emit(c.Tag, Error, message, fields...)
}

// F logs a message when the Level is set to Fatal or lower
func (c *Context) F(message string, fields ...interface{}) {
	c.Emitter.Emit(c.Tag, Fatal, message, fields...)
	_exit(1)
}
