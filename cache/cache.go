package cache

import (
	"log"
	"sync"
	"time"

	"github.com/kodernubie/gocommon/conf"
)

var listCreator map[string]ConnCreator = map[string]ConnCreator{}
var listConn map[string]Connection = map[string]Connection{}

var lock sync.RWMutex

func RegConnCreator(typeName string, creator ConnCreator) {

	lock.Lock()
	defer lock.Unlock()

	listCreator[typeName] = creator
}

func Conn(name ...string) Connection {

	targetName := "default"

	if len(name) > 0 {
		targetName = name[0]
	}

	lock.RLock()
	ret, exist := listConn[targetName]
	lock.RUnlock()

	if !exist {

		lock.Lock()
		defer lock.Unlock()

		// check wether other thread alredy set
		ret, exist = listConn[targetName]

		if !exist {

			typeName := conf.Str("CACHE_"+targetName+"_TYPE", "redis")
			creator, exist := listCreator[typeName]

			if !exist {
				log.Fatal("Unable to find cache connector connection ", typeName)
				return nil
			}

			var err error
			ret, err = creator(targetName)

			if err != nil {
				log.Fatal("Unable to create cache connection ", err.Error())
				return nil
			}

			listConn[targetName] = ret
		}
	}

	return ret
}

func Set(key, value string, expiration ...time.Duration) {

	Conn().Set(key, value, expiration...)
}

func SetNX(key, value string, expiration ...time.Duration) bool {

	return Conn().SetNX(key, value, expiration...)
}

func Get(key string) string {

	return Conn().Get(key)
}

func Incr(key string, expiration ...time.Duration) int {

	return Conn().Incr(key, expiration...)
}

func Decr(key string, expiration ...time.Duration) int {

	return Conn().Decr(key, expiration...)
}

func Del(key string) {

	Conn().Del(key)
}
