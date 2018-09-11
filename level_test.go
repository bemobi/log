package log_test

import (
	"testing"

	"github.com/bemobi/log"
)

func TestLevel(t *testing.T) {
	tests := []struct {
		Name  string
		Level string
		Want  log.Level
	}{
		{
			Name:  "Trace",
			Level: "trace",
			Want:  log.Trace,
		},
		{
			Name:  "Info",
			Level: "info",
			Want:  log.Info,
		},
		{
			Name:  "Error",
			Level: "error",
			Want:  log.Error,
		},
		{
			Name:  "Warn",
			Level: "warn",
			Want:  log.Warn,
		},
		{
			Name:  "Fatal",
			Level: "fatal",
			Want:  log.Fatal,
		},
		{
			Name:  "Invalid",
			Level: "",
			Want:  log.Info,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			log.SetLevel(test.Level)
			if log.Default.Level != test.Want {
				t.Errorf("got %d; want %d", log.Default.Level, test.Want)
			}
		})
	}
}
