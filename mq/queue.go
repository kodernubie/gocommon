package mq

import (
	"log"
	"sync"

	"github.com/kodernubie/gocommon/conf"
)

var listCreator map[string]ConnectionCreator = map[string]ConnectionCreator{}
var listConnection map[string]Connection = map[string]Connection{}

var lock sync.RWMutex

func RegConnectionCreator(typeName string, creator ConnectionCreator) {

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
	ret, exist := listConnection[targetName]
	lock.RUnlock()

	if !exist {

		lock.Lock()
		defer lock.Unlock()

		// check wether another thread set the value
		ret, exist = listConnection[targetName]

		if !exist {
			typeName := conf.Str("CACHE_"+targetName+"_TYPE", "rabbitmq")
			creator, exist := listCreator[typeName]

			if !exist {
				log.Fatal("Unable to find queue creator ", typeName)
				return nil
			}

			var err error
			ret, err = creator(typeName)

			if err != nil {
				log.Fatal("Unable to create queue connection ", err.Error())
				return nil
			}

			listConnection[targetName] = ret
		}
	}

	return ret
}

func Publish(queueName string, payload []byte) {

	Conn().Publish(queueName, payload)
}

func Subscribe(queueName string, handler ItemHandler) {

	Conn().Subscribe(queueName, handler)
}

func GroupSubscribe(queueName, groupName string, handler ItemHandler) {

	Conn().GroupSubscribe(queueName, groupName, handler)
}
