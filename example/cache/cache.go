package main

import (
	"time"

	"github.com/kodernubie/gocommon/cache"
	"github.com/kodernubie/gocommon/log"
)

func main() {

	log.Info("Cache Test")

	cache.Set("testx", "satu", 5*time.Minute)

	log.Info("Value =>", cache.Get("testx"))

	//test setnx, only set when value is not set

	if cache.SetNX("testx", "dua", 5*time.Minute) {
		log.Info("Success SetNX")
	} else {
		log.Info("Failed SetNX")
	}

	log.Info("After SetNX Value =>", cache.Get("testx"))

	cache.Incr("testCounter", 5*time.Minute)
	cache.Incr("testCounter", 5*time.Minute)
	cache.Incr("testCounter", 5*time.Minute)

	log.Info("counter =>", cache.Get("testCounter"))
}
