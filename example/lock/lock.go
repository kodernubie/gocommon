package main

import (
	"time"

	"github.com/kodernubie/gocommon/lock"
	"github.com/kodernubie/gocommon/log"
)

func main() {

	log.Info("Try to get lock....")

	lk := lock.Lock("test", 20*time.Second)
	defer lk.Unlock()

	log.Info("Lock acquired..")

	time.Sleep(10 * time.Second)

	log.Info("Exit")
}
