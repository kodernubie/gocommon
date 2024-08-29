package main

import (
	"github.com/kodernubie/gocommon/http"
	"github.com/kodernubie/gocommon/log"
)

func main() {

	log.Info("HTTP Server")

	http.Start()
}
