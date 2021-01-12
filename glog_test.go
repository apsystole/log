package log_test

import "github.com/apsystole/log"

func ExamplePrint() {
	log.Print("Test")
	// Output:
	// {"message":"Test","severity":"INFO"}
}

func ExampleNoticef() {
	log.Noticef("Hello %q!", "Google")
	// Output:
	// {"message":"Hello \"Google\"!","severity":"NOTICE"}
}

func ExampleWarningj() {
	log.Warningj("Warning", map[string]string{"component": "app"})
	// Output:
	// {"component":"app","message":"Warning","severity":"WARNING"}
}
