package log_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"
	"testing"

	"github.com/bemobi/log"
)

func TestLog(t *testing.T) {
	sink := &bytes.Buffer{}
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	tests := []struct {
		level   log.Level
		emitter rune
		name    string
		tag     string
		msg     string
		values  []interface{}
		want    map[string]interface{}
	}{
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Floats",
			tag:     "test",
			values:  []interface{}{"float32", float32(0.2), "float64", float64(0.271727127)},
			want: map[string]interface{}{
				"level":   "info",
				"msg":     "",
				"tag":     "test",
				"float32": 0.2,
				"float64": 0.271727127,
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Ints",
			tag:     "test",
			values:  []interface{}{"int", 1777, "int32", 127718, "int64", 100010},
			want: map[string]interface{}{
				"level": "info",
				"msg":   "",
				"tag":   "test",
				"int":   float64(1777),
				"int32": float64(127718),
				"int64": float64(100010),
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Uints",
			tag:     "test",
			values:  []interface{}{"uint", 1777, "uint32", 127718, "uint64", 100010},
			want: map[string]interface{}{
				"level":  "info",
				"msg":    "",
				"tag":    "test",
				"uint":   float64(uint(1777)),
				"uint32": float64(uint32(127718)),
				"uint64": float64(uint64(100010)),
			},
		},
		{
			level:   log.Error,
			emitter: 'T',
			name:    "T",
			tag:     "test",
			values:  []interface{}{"uint", 1777, "uint32", 127718, "uint64", 100010},
			want:    map[string]interface{}{},
		},
		{
			level:   log.Error,
			emitter: 'I',
			name:    "I",
			tag:     "test",
			values:  []interface{}{"uint", 1777, "uint32", 127718, "uint64", 100010},
			want:    map[string]interface{}{},
		},
		{
			level:   log.Error,
			emitter: 'W',
			name:    "W",
			tag:     "test",
			values:  []interface{}{"uint", 1777, "uint32", 127718, "uint64", 100010},
			want:    map[string]interface{}{},
		},
		{
			level:   log.Error,
			emitter: 'E',
			name:    "E",
			tag:     "test",
			values: []interface{}{
				"uint", 1777,
				"uint32", 127718,
				"uint64", 100010,
				"error", errors.New("Hello World"),
				"t", newT(),
			},
			want: map[string]interface{}{
				"level":  "error",
				"msg":    "",
				"tag":    "test",
				"uint":   float64(uint(1777)),
				"uint32": float64(uint32(127718)),
				"uint64": float64(uint64(100010)),
				"error":  "Hello World",
				"t":      "Hi - Ho",
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Quotes",
			tag:     "test",
			msg:     "xpto",
			values: []interface{}{
				"a", `"123"`,
				"b", `"x":"y"`,
				"c", errors.New(`"error"`),
				"d", []byte(`uno,"2",'three'`),
				"e", [3]int{1, 2, 3},
				"f", map[string]bool{
					"1": true,
				},
			},
			want: map[string]interface{}{
				"level": "info",
				"msg":   "xpto",
				"tag":   "test",
				"a":     `"123"`,
				"b":     `"x":"y"`,
				"c":     `"error"`,
				"d":     `uno,"2",'three'`,
				"e":     `[1 2 3]`,
				"f":     `map[1:true]`,
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Line Breaks and Tabs",
			tag:     "test",
			msg:     "a\nb\t\tc",
			values: []interface{}{
				"x", "\t\t\ty\n\n\n",
			},
			want: map[string]interface{}{
				"level": "info",
				"msg":   "a\nb\t\tc",
				"tag":   "test",
				"x":     "\t\t\ty\n\n\n",
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Bad Escaping",
			tag:     "test",
			msg:     `\t\h\i\s\ \s\h\o\u\l\d\ \n\o\t\ \b\r\e\a\k`,
			values: []interface{}{
				"x", `\a\n\d\ \a\l\s\o\ \t\h\i\s\ \o\n\e`,
			},
			want: map[string]interface{}{
				"level": "info",
				"msg":   `\t\h\i\s\ \s\h\o\u\l\d\ \n\o\t\ \b\r\e\a\k`,
				"tag":   "test",
				"x":     `\a\n\d\ \a\l\s\o\ \t\h\i\s\ \o\n\e`,
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Some other cases, with UTF8",
			tag:     "test",
			msg:     `Olá usuário! Esse é um teþte!☺`,
			values: []interface{}{
				"x", `Olá usuário! Esse é um teþte2!☺`,
			},
			want: map[string]interface{}{
				"level": "info",
				"msg":   `Olá usuário! Esse é um teþte!☺`,
				"tag":   "test",
				"x":     `Olá usuário! Esse é um teþte2!☺`,
			},
		},
		{
			level:   log.Info,
			emitter: 'I',
			name:    "Japanese Handling",
			tag:     "test",
			msg:     `新 あたら しい 記事 きじ を 書 か こうという`,
			values: []interface{}{
				"x", `ど、 息子 むすこ を 産 う んだ 後`,
			},
			want: map[string]interface{}{
				"level": "info",
				"msg":   `新 あたら しい 記事 きじ を 書 か こうという`,
				"tag":   "test",
				"x":     `ど、 息子 むすこ を 産 う んだ 後`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sink.Reset()
			log.Default.Level = test.level
			switch test.emitter {
			case 'T':
				log.T(test.tag, test.msg, test.values...)
			case 'I':
				log.I(test.tag, test.msg, test.values...)
			case 'W':
				log.W(test.tag, test.msg, test.values...)
			case 'E':
				log.E(test.tag, test.msg, test.values...)
			case 'F':
				log.F(test.tag, test.msg, test.values...)
			}

			got := make(map[string]interface{})
			if bytes := sink.Bytes(); len(bytes) > 0 {
				if err := json.Unmarshal(bytes, &got); err != nil {
					t.Fatalf("invalid json logging: %s", string(bytes))
				}
			}

			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("\nwant: %v\ngot: %v", test.want, got)
			}
		})
	}
}

//go:noinline
func noOpHook(log.Level, []byte) {}

func BenchmarkEmitSingleFieldString(b *testing.B) { setupBenchmark(1, true, nil, b) }
func BenchmarkEmitFourFieldsString(b *testing.B)  { setupBenchmark(4, true, nil, b) }
func BenchmarkEmitTenFieldsString(b *testing.B)   { setupBenchmark(10, true, nil, b) }

func BenchmarkEmitSingleFieldInt(b *testing.B) { setupBenchmark(1, false, nil, b) }
func BenchmarkEmitFourFieldsInt(b *testing.B)  { setupBenchmark(4, false, nil, b) }
func BenchmarkEmitTenFieldsInt(b *testing.B)   { setupBenchmark(10, false, nil, b) }

func BenchmarkEmitWithHookSingleFieldString(b *testing.B) { setupBenchmark(1, true, noOpHook, b) }
func BenchmarkEmitWithHookFourFieldsString(b *testing.B)  { setupBenchmark(4, true, noOpHook, b) }
func BenchmarkEmitWithHookTenFieldsString(b *testing.B)   { setupBenchmark(10, true, noOpHook, b) }

func BenchmarkEmitWithHookSingleFieldInt(b *testing.B) { setupBenchmark(1, false, noOpHook, b) }
func BenchmarkEmitWithHookFourFieldsInt(b *testing.B)  { setupBenchmark(4, false, noOpHook, b) }
func BenchmarkEmitWithHookTenFieldsInt(b *testing.B)   { setupBenchmark(10, false, noOpHook, b) }

func setupBenchmark(howManyFields int, asString bool, hook log.Hook, b *testing.B) {
	log.SetHook(hook)
	defer log.SetHook(nil)

	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	fs := make([]interface{}, 0)
	for i := 0; i < howManyFields; i++ {
		s := interface{}(i)
		if asString {
			s = string(i)
		}
		fs = append(fs, fmt.Sprintf("field_%d", i), fmt.Sprintf(`\t%s\n\\\`, s))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.I("test", "hello", fs...)
	}
	b.StopTimer()
	b.ReportAllocs()
}

type Test struct {
	A string
	B string
}

func newT() *Test { return &Test{"Hi", "Ho"} }

func (tst *Test) String() string {
	return fmt.Sprintf("%s - %s", tst.A, tst.B)
}

func BenchmarkLogStringer(b *testing.B) {
	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	tst := newT()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.I("test", "hello", "tst", tst)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkLogError(b *testing.B) {
	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.I("test", "hello", "err", errors.New("hello world"))
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkContext10Fields(b *testing.B) {
	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	logger := log.C("test",
		"one", "one",
		"two", "two",
		"three", "three",
		"four", "four",
		"five", "five",
		"six", "six",
		"seven", "seven",
		"eight", "eight",
		"nine", "nine",
		"ten", "ten",
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.I("world")
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkContext10FieldsConcurrent(b *testing.B) {
	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	logger := log.C("test",
		"one", "one",
		"two", "two",
		"three", "three",
		"four", "four",
		"five", "five",
		"six", "six",
		"seven", "seven",
		"eight", "eight",
		"nine", "nine",
		"ten", "ten",
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.I("world")
		}
	})
	b.StopTimer()
	b.ReportAllocs()
}

func TestConcurrentWrites(t *testing.T) {
	sink := ioutil.Discard
	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	err := errors.New("hello world")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			log.I("test", "hello1", "err", err)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			log.I("test", "hello2", "err", err)
		}
	}()
	wg.Wait()
}

func TestHook(t *testing.T) {
	rec := make(map[log.Level][]string)

	sink := &bytes.Buffer{}
	log.SetLevel("trace")

	log.SetHook(func(l log.Level, doc []byte) {
		rec[l] = append(rec[l], string(doc))
	})
	defer log.SetHook(nil)

	log.SetTestMode(true, sink)
	defer log.SetTestMode(false)

	const tag = "test"
	log.T(tag, "1", "a", 1)
	log.I(tag, "21")
	log.I(tag, "22")

	logger := log.C(tag, "x", 5, "y", 4)
	logger.W("3")
	logger.E("4", "err", errors.New("new error"))

	logger = logger.C("z", 3)
	logger.F("5")

	t.Run("Trace", func(t *testing.T) {
		if len(rec[log.Trace]) != 1 {
			t.Fatal("trace hook was not called once")
		}
		want := `{"tag":"test","level":"debug","msg":"1","a":1}` + "\n"
		got := rec[log.Trace][0]
		if got != want {
			t.Errorf("trace is wrong: want[%s] got [%s]", want, got)
		}
	})

	t.Run("Info", func(t *testing.T) {
		if len(rec[log.Info]) != 2 {
			t.Fatal("info hook was not called twice")
		}

		want := `{"tag":"test","level":"info","msg":"21"}` + "\n"
		got := rec[log.Info][0]
		if got != want {
			t.Errorf("info is wrong: want[%s] got [%s]", want, got)
		}

		want = `{"tag":"test","level":"info","msg":"22"}` + "\n"
		got = rec[log.Info][1]
		if got != want {
			t.Errorf("info is wrong: want[%s] got [%s]", want, got)
		}
	})

	t.Run("Warning", func(t *testing.T) {
		if len(rec[log.Warn]) != 1 {
			t.Fatal("warn hook was not called once")
		}

		want := `{"tag":"test","level":"warn","msg":"3","x":5,"y":4}` + "\n"
		got := rec[log.Warn][0]
		if got != want {
			t.Errorf("warn is wrong: want[%s] got [%s]", want, got)
		}
	})

	t.Run("Error", func(t *testing.T) {
		if len(rec[log.Error]) != 1 {
			t.Fatal("error hook was not called once")
		}

		want := `{"tag":"test","level":"error","msg":"4","x":5,"y":4,"err":"new error"}` + "\n"
		got := rec[log.Error][0]
		if got != want {
			t.Errorf("error is wrong: want[%s] got [%s]", want, got)
		}
	})

	t.Run("Fatal", func(t *testing.T) {
		if len(rec[log.Fatal]) != 1 {
			t.Fatal("fatal hook was not called once")
		}

		want := `{"tag":"test","level":"fatal","msg":"5","x":5,"y":4,"z":3}` + "\n"
		got := rec[log.Fatal][0]
		if got != want {
			t.Errorf("fatal is wrong: want[%s] got [%s]", want, got)
		}
	})
}
