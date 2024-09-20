package main

import (
	"github.com/kodernubie/gocommon/log"
	"github.com/kodernubie/gocommon/mq"
)

func main() {

	mq.GroupSubscribe("test", "listener", func(queueName string, payload []byte) {

		log.Info("Received :", string(payload))
	})

	var forever chan struct{}

	log.Info("Press ctrl+c to exit")
	<-forever
}
