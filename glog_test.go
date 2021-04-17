package log

// Test the public interface, but inject some private dependencies; run it from the same package.

import (
	"testing"
)

func TestPanic(t *testing.T) {
	defer func() {
		want := "a"
		if r := recover(); r != want {
			t.Errorf("got panic %q expected %q", r, want)
		}
	}()
	Panic("a")
}
