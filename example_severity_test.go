package log_test

import "github.com/apsystole/log"

// Example shows typical severity recommendations.
func Example() {
	log.Debug("the lowest level")
	log.Print("phase 2 processing ended") // or log.Info()
	log.Notice("shutting down normally...")
	log.Warning("difficulties, retrying shutdown...")
}

// Errors are used sparingly to draw attention of a human, sooner or later.
func Errors() {
	log.Error("a minor problem that a real human should act on")
	log.Critical("a major problem")
	log.Fatal("a major problem - aborting")
	log.Panic("a major problem - dumping stacktraces and aborting")
	log.Alert("waky waky dear on-duty person, eggs and baky")
	log.Emergency("core functionality is down and needs a human to rescue it")
}
