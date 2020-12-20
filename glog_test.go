package log_test

import "github.com/apsystole/log"

func ExamplePrint() {
	log.Print("Test")
	// Output:
	// {"message":"Test"}
}

func ExampleInfof() {
	log.Infof("Hello %q!", "Google")
	// Output:
	// {"message":"Hello \"Google\"!","severity":"INFO"}
}

func ExampleWarningj() {
	log.Warningj("Warning", map[string]string{"component": "app"})
	// Output:
	// {"component":"app","message":"Warning","severity":"WARNING"}
}
