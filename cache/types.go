package cache

import "time"

type Connection interface {
	Set(key, value string, expiration ...time.Duration)
	SetNX(key, value string, expiration ...time.Duration) bool
	Get(key string) string
	Incr(key string, expiration ...time.Duration) int
	Decr(key string, expiration ...time.Duration) int
	Del(key string)
}

type ConnCreator func(configName string) (Connection, error)
