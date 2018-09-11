package log_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bemobi/log"
)

func TestContext(t *testing.T) {
	sink := &bytes.Buffer{}
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	type Emit struct {
		Level   log.Level
		Message string
		Fields  []interface{}
	}

	for _, tc := range []struct {
		Name     string
		Context  *log.Context
		Level    log.Level
		Emit     []Emit
		Want     []map[string]interface{}
		ExitCode int
	}{
		{
			Name:    "Without Fields",
			Context: log.C("TAG"),
			Level:   log.Info,
			Emit: []Emit{
				{Level: log.Warn, Message: "m1"},
				{Level: log.Error, Message: "m2"},
			},
			Want: []map[string]interface{}{
				{"tag": "TAG", "level": "warn", "msg": "m1"},
				{"tag": "TAG", "level": "error", "msg": "m2"},
			},
		},
		{
			Name:    "With Two Fields",
			Context: log.C("TAG", "one", "first", "two", "second"),
			Level:   log.Info,
			Emit: []Emit{
				{Level: log.Warn, Message: "m1"},
				{Level: log.Error, Message: "m2"},
			},
			Want: []map[string]interface{}{
				{"tag": "TAG", "one": "first", "two": "second", "level": "warn", "msg": "m1"},
				{"tag": "TAG", "one": "first", "two": "second", "level": "error", "msg": "m2"},
			},
		},
		{
			Name:    "With Lower Level",
			Context: log.C("TAG", "one", "first", "two", "second"),
			Level:   log.Fatal,
			Emit: []Emit{
				{Level: log.Warn, Message: "m1"},
				{Level: log.Error, Message: "m2"},
			},
			Want: []map[string]interface{}{},
		},
		{
			Name:    "Nested",
			Context: log.C("TAG", "one", "first").C("two", "second"),
			Level:   log.Trace,
			Emit: []Emit{
				{Level: log.Trace, Message: "m1"},
				{Level: log.Info, Message: "m2"},
				{Level: log.Fatal, Message: "m3"},
			},
			Want: []map[string]interface{}{
				{"tag": "TAG", "one": "first", "two": "second", "level": "debug", "msg": "m1"},
				{"tag": "TAG", "one": "first", "two": "second", "level": "info", "msg": "m2"},
				{"tag": "TAG", "one": "first", "two": "second", "level": "fatal", "msg": "m3"},
			},
			ExitCode: 1,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Context.Emitter.Level = tc.Level

			for _, emit := range tc.Emit {
				switch emit.Level {
				case log.Trace:
					tc.Context.T(emit.Message, emit.Fields...)
				case log.Info:
					tc.Context.I(emit.Message, emit.Fields...)
				case log.Warn:
					tc.Context.W(emit.Message, emit.Fields...)
				case log.Error:
					tc.Context.E(emit.Message, emit.Fields...)
				case log.Fatal:
					tc.Context.F(emit.Message, emit.Fields...)
				}
			}

			scanner := bufio.NewScanner(sink)
			index := 0
			for scanner.Scan() {
				got := make(map[string]interface{})
				err := json.Unmarshal(scanner.Bytes(), &got)
				if err != nil {
					t.Fatalf("invalid json format: %v", err)
				}
				want := tc.Want[index]
				if !reflect.DeepEqual(want, got) {
					t.Fatalf("invalid output:\nwant: %v\ngot: %v", want, got)
				}
				index++
			}
			if index != len(tc.Want) {
				t.Fatalf("invalid output count:\nwant: %v\ngot: %v", len(tc.Want), index)
			}

			if log.LastExitCode() != tc.ExitCode {
				t.Fatalf("invalid exit code:\nwant: %v\ngot: %v", tc.ExitCode, log.LastExitCode())
			}
		})
	}

}
