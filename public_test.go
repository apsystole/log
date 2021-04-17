package log_test

// Test using only public interface, so run it from a separate package.

import (
	"io"
	"testing"

	"github.com/apsystole/log"
)

func ExampleNotice() {
	log.Notice("hello", 1, "!")
	log.Noticeln("hello", 2, "!")
	// Output:
	// {"message":"hello1!","severity":"NOTICE"}
	// {"message":"hello 2 !\n","severity":"NOTICE"}
}

func ExampleNoticeln() {
	log.Notice("hello", 1, "!")
	log.Noticeln("hello", 2, "!")
	// Output:
	// {"message":"hello1!","severity":"NOTICE"}
	// {"message":"hello 2 !\n","severity":"NOTICE"}
}

func ExampleNoticef() {
	err := io.EOF
	name := "my blog.txt"
	log.Noticef("while reading file %q ignoring: %v", name, err)
	// Output:
	// {"message":"while reading file \"my blog.txt\" ignoring: EOF","severity":"NOTICE"}
}

func ExampleNoticej() {
	obj := struct {
		Seq       int
		Component string
	}{
		Seq:       42,
		Component: "app",
	}
	log.Noticej("warning", obj)
	// Output:
	// {"Component":"app","Seq":42,"message":"warning","severity":"NOTICE"}
}

func BenchmarkDebug(b *testing.B) {
	l := log.New(nil, "", 0)
	for i := 0; i < b.N; i++ {
		l.Debug("test")
	}
}
