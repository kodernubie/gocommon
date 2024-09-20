package main

import (
	"strconv"
	"time"

	"github.com/kodernubie/gocommon/log"
	"github.com/kodernubie/gocommon/mq"
)

func main() {

	for i := 0; i < 10; i++ {

		data := []byte(strconv.Itoa(i) + " =>" + time.Now().Format(time.RFC3339))

		log.Info("Publish :", string(data))
		mq.Publish("test", data)
		time.Sleep(time.Second)
	}

	log.Info("Selesai")
}
