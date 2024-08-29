package db

import (
	"log"
	"sync"

	"github.com/kodernubie/gocommon/conf"
)

var listCreator map[string]ConnectionCreator = map[string]ConnectionCreator{}
var listConn map[string]Connection = map[string]Connection{}

var lock sync.RWMutex

func RegConnCreator(typeName string, creator ConnectionCreator) {

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

		// check if another thread already set
		ret, exist = listConn[targetName]

		if !exist {

			typeName := conf.Str("DB_"+targetName+"_TYPE", "mongo")
			creator, exist := listCreator[typeName]

			if !exist {
				log.Fatal("Unable to find DB creator ", typeName)
				return nil
			}

			ret, err := creator(typeName)

			if err != nil {
				log.Fatal("Unable to create DB connection ", err)
				return nil
			}

			listConn[targetName] = ret
		}
	}

	return ret
}

func Insert(obj interface{}) error {

	return nil
}

func Update(obj interface{}) error {

	return nil
}

func Delete(obj interface{}) error {

	return nil
}

func Find(obj interface{}) ([]interface{}, error) {

	return nil, nil
}

func Raw() error {

	return nil
}
