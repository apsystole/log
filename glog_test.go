package log

// Test the public interface, but inject some private dependencies; run it from the same package.

import (
	"bytes"
	"testing"
)

func TestPanic(t *testing.T) {
	// Arrange
	wantJson := "{\"message\":\"a\",\"severity\":\"CRITICAL\"}\n"
	wantPanic := "a"
	buf := &bytes.Buffer{}
	l := New(buf, "", 0)

	// Assert
	defer func() {
		if gotPanic := recover(); gotPanic != wantPanic {
			t.Errorf("unexpected panic, got:\n%q\nexpected:\n%q\n", gotPanic, wantPanic)
		}
		if wantJson != buf.String() {
			t.Errorf("unexpected output, got:\n%q\nexpected:\n%q\n", buf.String(), wantJson)
		}
	}()

	// Act
	l.Panic("a")
}

func BenchmarkDebug(b *testing.B) {
	buf := &bytes.Buffer{} // quite unrealistic, as a write to file here is about 14 000 ns
	l := New(buf, "", 0)
	for i := 0; i < b.N; i++ {
		l.Debug("test")
		buf.Reset()
	}
}
