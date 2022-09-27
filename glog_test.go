package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestLogger_Panic(t *testing.T) {
	// Arrange
	wantJSON := `{"message":"a","severity":"CRITICAL"}
`
	wantPanic := "a"
	buf := &bytes.Buffer{}
	l := New(buf, "", 0)

	// Assert
	defer func() {
		if gotPanic := recover(); gotPanic != wantPanic {
			t.Errorf("unexpected panic, got:\n%q\nexpected:\n%q\n", gotPanic, wantPanic)
		}
		if wantJSON != buf.String() {
			t.Errorf("unexpected output, got:\n%q\nexpected:\n%q\n", buf.String(), wantJSON)
		}
	}()

	// Act
	l.Panic("a")
}

func TestLogger_Print(t *testing.T) {
	type fields struct {
		trace []byte
	}

	tests := []struct {
		name   string
		fields fields
		arg    string
		want   string
	}{{
		name: "empty message",
		arg:  "",
		want: `{"message":"","severity":"INFO"}
`,
	}, {
		name: "ending newline",
		arg:  "test\n",
		want: `{"message":"test\n","severity":"INFO"}
`,
	}, {
		name: "non-ending newline",
		arg:  "line\ntest",
		want: `{"message":"line\ntest","severity":"INFO"}
`,
	}, {
		name: "ampersand",
		arg:  "m&m",
		want: `{"message":"m&m","severity":"INFO"}
`,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// AAA: Arrange
			buf := &bytes.Buffer{}

			l := &Logger{
				out:   buf,
				err:   nil,
				trace: tt.fields.trace,
			}

			// AAA: Act
			l.Print(tt.arg)

			// AAA: Assert
			if tt.want != buf.String() {
				t.Errorf("unexpected output, got:\n%q\nexpected:\n%q\n", buf.String(), tt.want)
			}
			if !json.Valid(buf.Bytes()) {
				t.Errorf("output is not a valid JSON:\n%q\n", buf.Bytes())
			}
		})
	}
}

func TestLogger_Debugj(t *testing.T) {
	buf := &bytes.Buffer{}

	type fields struct {
		out   io.Writer
		err   io.Writer
		trace []byte
	}

	type args struct {
		msg string
		v   interface{}
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{{
		name:   "easy struct",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v:   struct{ Text string }{Text: "t"},
		},
		want: `{"message":"test","severity":"DEBUG","Text":"t"}
`,
	}, {
		name: "tracing the numeral",
		fields: fields{
			out:   buf,
			trace: []byte(`"123"`),
		},
		args: args{
			msg: "test",
			v:   struct{ Text string }{Text: "t"},
		},
		want: `{"message":"test","severity":"DEBUG","logging.googleapis.com/trace":"123","Text":"t"}
`,
	}, {
		name:   "empty struct",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v:   struct{}{},
		},
		want: `{"message":"test","severity":"DEBUG"}
`,
	}, {
		name:   "typed nil",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v:   (*string)(nil),
		},
		want: `{"message":"test","severity":"DEBUG","value":null}
`,
	}, {
		name:   "untyped nil",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v:   nil,
		},
		want: `{"message":"test","severity":"DEBUG","value":null}
`,
	}, {
		name:   "ampersand in a key",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v: struct {
				MandM string `json:"m&m"`
			}{
				MandM: "not brown",
			},
		},
		want: `{"message":"test","severity":"DEBUG","m&m":"not brown"}
`,
	}, {
		name:   "ampersand in a string value",
		fields: fields{out: buf},
		args: args{
			msg: "test",
			v: struct {
				MandM string
			}{
				MandM: "brown&banned",
			},
		},
		want: `{"message":"test","severity":"DEBUG","MandM":"brown&banned"}
`,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				out:   tt.fields.out,
				err:   tt.fields.err,
				trace: tt.fields.trace,
			}
			l.Debugj(tt.args.msg, tt.args.v)
			if tt.want != buf.String() {
				t.Errorf("unexpected output, got:\n%q\nexpected:\n%q\n", buf.String(), tt.want)
			}
			if !json.Valid(buf.Bytes()) {
				t.Errorf("output is not a valid JSON:\n%q\n", buf.Bytes())
			}
			buf.Reset()
		})
	}
}

// failingType fakes an error every time it's marshaled as JSON.
type failingType int

func (t failingType) MarshalJSON() ([]byte, error) {
	return nil, errors.New("e")
}

func TestPrintj_FailingMarshaler(t *testing.T) {
	// Arrange
	wantJSON := `{"message":"a","severity":"INFO","logLibMsg":"cannot marshal the argument as jsonPayload"}
`
	buf := &bytes.Buffer{}
	l := New(buf, "", 0)
	var ft failingType

	// Act
	l.Printj("a", ft)

	// Assert
	if wantJSON != buf.String() {
		t.Errorf("unexpected output, got:\n%q\nexpected:\n%q\n", buf.String(), wantJSON)
	}
}

// panickingType fakes a panic every time it's marshaled as JSON.
type panickingType int

func (t panickingType) MarshalJSON() ([]byte, error) {
	panic("p")
}

func TestPrintj_PanickingMarshaler(t *testing.T) {
	// Arrange
	wantPanic := "p"
	var pt panickingType

	// Assert
	defer func() {
		if gotPanic := recover(); gotPanic != wantPanic {
			t.Errorf("unexpected panic, got:\n%q\nexpected:\n%q\n", gotPanic, wantPanic)
		}
	}()

	// Act
	Printj("a", pt)
}

func BenchmarkDebugf(b *testing.B) {
	buf := &bytes.Buffer{} // quite unrealistic, as a write to file here is about 10_000 ns
	l := New(buf, "", 0)
	for i := 0; i < b.N; i++ {
		l.Debugf("%q", "test")
		buf.Reset()
	}
}

func BenchmarkDebugjHundred(b *testing.B) {
	buf := &bytes.Buffer{}
	msg := &struct {
		Field00 string
		Field01 string
		Field02 string
		Field03 string
		Field04 string
		Field05 string
		Field06 string
		Field07 string
		Field08 string
		Field09 string
		Field10 string
		Field11 string
		Field12 string
		Field13 string
		Field14 string
		Field15 string
		Field16 string
		Field17 string
		Field18 string
		Field19 string
		Field20 string
		Field21 string
		Field22 string
		Field23 string
		Field24 string
		Field25 string
		Field26 string
		Field27 string
		Field28 string
		Field29 string
		Field30 string
		Field31 string
		Field32 string
		Field33 string
		Field34 string
		Field35 string
		Field36 string
		Field37 string
		Field38 string
		Field39 string
		Field40 string
		Field41 string
		Field42 string
		Field43 string
		Field44 string
		Field45 string
		Field46 string
		Field47 string
		Field48 string
		Field49 string
		Field50 string
		Field51 string
		Field52 string
		Field53 string
		Field54 string
		Field55 string
		Field56 string
		Field57 string
		Field58 string
		Field59 string
		Field60 string
		Field61 string
		Field62 string
		Field63 string
		Field64 string
		Field65 string
		Field66 string
		Field67 string
		Field68 string
		Field69 string
		Field70 string
		Field71 string
		Field72 string
		Field73 string
		Field74 string
		Field75 string
		Field76 string
		Field77 string
		Field78 string
		Field79 string
		Field80 string
		Field81 string
		Field82 string
		Field83 string
		Field84 string
		Field85 string
		Field86 string
		Field87 string
		Field88 string
		Field89 string
		Field90 string
		Field91 string
		Field92 string
		Field93 string
		Field94 string
		Field95 string
		Field96 string
		Field97 string
		Field98 string
		Field99 string
	}{
		Field00: "test",
		Field01: "test",
		Field02: "test",
		Field03: "test",
		Field04: "test",
		Field05: "test",
		Field06: "test",
		Field07: "test",
		Field08: "test",
		Field09: "test",
		Field10: "test",
		Field11: "test",
		Field12: "test",
		Field13: "test",
		Field14: "test",
		Field15: "test",
		Field16: "test",
		Field17: "test",
		Field18: "test",
		Field19: "test",
		Field20: "test",
		Field21: "test",
		Field22: "test",
		Field23: "test",
		Field24: "test",
		Field25: "test",
		Field26: "test",
		Field27: "test",
		Field28: "test",
		Field29: "test",
		Field30: "test",
		Field31: "test",
		Field32: "test",
		Field33: "test",
		Field34: "test",
		Field35: "test",
		Field36: "test",
		Field37: "test",
		Field38: "test",
		Field39: "test",
		Field40: "test",
		Field41: "test",
		Field42: "test",
		Field43: "test",
		Field44: "test",
		Field45: "test",
		Field46: "test",
		Field47: "test",
		Field48: "test",
		Field49: "test",
		Field50: "test",
		Field51: "test",
		Field52: "test",
		Field53: "test",
		Field54: "test",
		Field55: "test",
		Field56: "test",
		Field57: "test",
		Field58: "test",
		Field59: "test",
		Field60: "test",
		Field61: "test",
		Field62: "test",
		Field63: "test",
		Field64: "test",
		Field65: "test",
		Field66: "test",
		Field67: "test",
		Field68: "test",
		Field69: "test",
		Field70: "test",
		Field71: "test",
		Field72: "test",
		Field73: "test",
		Field74: "test",
		Field75: "test",
		Field76: "test",
		Field77: "test",
		Field78: "test",
		Field79: "test",
		Field80: "test",
		Field81: "test",
		Field82: "test",
		Field83: "test",
		Field84: "test",
		Field85: "test",
		Field86: "test",
		Field87: "test",
		Field88: "test",
		Field89: "test",
		Field90: "test",
		Field91: "test",
		Field92: "test",
		Field93: "test",
		Field94: "test",
		Field95: "test",
		Field96: "test",
		Field97: "test",
		Field98: "test",
		Field99: "test",
	}
	l := New(buf, "", 0)
	for i := 0; i < b.N; i++ {
		l.Debugj("test", msg)
		buf.Reset()
	}
}

func BenchmarkDebugjTen(b *testing.B) {
	buf := &bytes.Buffer{}
	msg := &struct {
		Field00 string
		Field01 string
		Field02 string
		Field03 string
		Field04 string
		Field05 string
		Field06 string
		Field07 string
		Field08 string
		Field09 string
		Field10 string
	}{
		Field00: "test",
		Field01: "test",
		Field02: "test",
		Field03: "test",
		Field04: "test",
		Field05: "test",
		Field06: "test",
		Field07: "test",
		Field08: "test",
		Field09: "test",
		Field10: "test",
	}
	l := New(buf, "", 0)
	for i := 0; i < b.N; i++ {
		l.Debugj("test", msg)
		buf.Reset()
	}
}

func BenchmarkJsonTen(b *testing.B) {
	buf := &bytes.Buffer{}
	msg := &struct {
		Field00 string
		Field01 string
		Field02 string
		Field03 string
		Field04 string
		Field05 string
		Field06 string
		Field07 string
		Field08 string
		Field09 string
		Field10 string
	}{
		Field00: "test",
		Field01: "test",
		Field02: "test",
		Field03: "test",
		Field04: "test",
		Field05: "test",
		Field06: "test",
		Field07: "test",
		Field08: "test",
		Field09: "test",
		Field10: "test",
	}
	for i := 0; i < b.N; i++ {
		_ = json.NewEncoder(buf).Encode(msg)
		_ = json.Unmarshal(buf.Bytes(), msg)
		buf.Reset()
	}
}

func BenchmarkStdlibTen(b *testing.B) {
	buf := &bytes.Buffer{}
	l := New(buf, "", 0)
	msg := &struct {
		Field00 string
		Field01 string
		Field02 string
		Field03 string
		Field04 string
		Field05 string
		Field06 string
		Field07 string
		Field08 string
		Field09 string
		Field10 string
	}{
		Field00: "test",
		Field01: "test",
		Field02: "test",
		Field03: "test",
		Field04: "test",
		Field05: "test",
		Field06: "test",
		Field07: "test",
		Field08: "test",
		Field09: "test",
		Field10: "test",
	}
	for i := 0; i < b.N; i++ {
		logjStdlib(debugsev, l, "test", msg)
		buf.Reset()
	}
}

// logjStdlib is only here to benchmark the stdlib "encoding/json"
// approach. Hopefully our method is much faster than stdlib.
func logjStdlib(s severity, l *Logger, msg string, j interface{}) {
	entry := make(map[string]json.RawMessage)

	if buf, err := json.Marshal(j); err != nil {
		panic(err)
	} else if err := json.Unmarshal(buf, &entry); err != nil {
		panic(err)
	}

	if v := msg; v != "" {
		entry["message"], _ = json.Marshal(v)
	}
	if v := l.trace; len(v) != 0 {
		entry["logging.googleapis.com/trace"] = v
	}
	entry["severity"], _ = s.MarshalJSON()

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = json.NewEncoder(l.writer(s)).Encode(entry)
}

func TestForRequest(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name      string
		projectID string
		args      args
		want      *Logger
	}{{
		name: "no tracing header",
		args: args{req: &http.Request{Header: http.Header{}}},
		want: &Logger{},
	}, {
		name: "empty tracing header",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{""},
		}}},
		want: &Logger{},
	}, {
		name: "tracing header without project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"00000000000000000000000000000001/1;"},
		}}},
		want: &Logger{},
	}, {
		name:      "basic tracing",
		projectID: "my-project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"00000000000000000000000000000001/1;o=1"},
		}}},
		want: &Logger{
			trace: []byte(`"projects/my-project/traces/00000000000000000000000000000001"`),
		},
	}, {
		name:      "tracing header without the o option",
		projectID: "my-project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"00000000000000000000000000000001/1"},
		}}},
		want: &Logger{
			trace: []byte(`"projects/my-project/traces/00000000000000000000000000000001"`),
		},
	}, {
		name:      "o=0 header disables tracing",
		projectID: "my-project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"00000000000000000000000000000001/1;o=0"},
		}}},
		want: &Logger{},
	}, {
		name:      "bad header no tid",
		projectID: "my-project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"/123;o=1"},
		}}},
		want: &Logger{},
	}, {
		name:      "bad header malformed tid",
		projectID: "my-project",
		args: args{req: &http.Request{Header: http.Header{
			"X-Cloud-Trace-Context": []string{"&/123;o=1"},
		}}},
		want: &Logger{},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ProjectID = tt.projectID

			got := ForRequest(tt.args.req)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
