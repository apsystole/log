package log_test

import "github.com/apsystole/log"

func ExampleSevRecommendations() {
	log.Debug("the lowest level")
	log.Print("phase 2 processing ended") // or Info()
	log.Notice("shutting down normally...")
	log.Warning("retrying shutdown...")
}

func Errors() {
	// Use these sparingly to draw attention of a human (sooner or later):
	log.Error("a minor problem")
	log.Critical("a major problem")
	log.Fatal("a major problem - aborting")
	log.Panic("a major problem - dumping stacktraces and aborting")
	log.Alert("intervetion of a real human likely required")
	log.Emergency("waky waky dear on-duty person, eggs and baky")
}
