package main

import (
	golog "log"

	"github.com/kodernubie/gocommon/log"
)

func main() {

	log.Info("Log Example")

	log.Error("This is error log")
	log.Errorf("This is error with parameter %d", 1)

	log.Debug("This is debug log")
	log.Debugf("This is bebug with parameter --%d===%s--", 579, "clear")

	golog.Println("info.....")
}
