package http

import (
	"strings"
	"sync"

	"github.com/kodernubie/gocommon/conf"
	"github.com/kodernubie/gocommon/log"
)

var creators map[string]Creator = map[string]Creator{}
var svrs map[string]Server = map[string]Server{}

var mtx sync.RWMutex

func RegCreator(typeName string, creator Creator) {
	creators[typeName] = creator
}

func Svr(name ...string) Server {

	targetName := "default"

	if len(name) > 0 {
		targetName = name[0]
	}

	mtx.RLock()
	ret, exist := svrs[targetName]
	mtx.RUnlock()

	if !exist {

		mtx.Lock()
		defer mtx.Unlock()

		ret, exist = svrs[targetName]

		if !exist {

			typeName := conf.Str("HTTP_" + strings.ToUpper(targetName) + "_TYPE")

			if typeName == "" {
				typeName = "fiber" // default to fiber
			}

			creator, exist := creators[typeName]

			if !exist {
				log.Fatal("Unable to create http server because creaator is not exist ", targetName,
					conf.Str("HTTP_"+strings.ToUpper(targetName)+"_TYPE"))
				return nil
			}

			var err error
			ret, err = creator(targetName)

			if err != nil {
				log.Fatal("Unable to create http server", targetName, err)
				return nil
			}

			svrs[targetName] = ret
		}
	}

	return ret
}

func Get(url string, handlers ...Handler) {

	Svr().Get(url, handlers...)
}

func Post(url string, handlers ...Handler) {

	Svr().Post(url, handlers...)
}

func Put(url string, handlers ...Handler) {

	Svr().Put(url, handlers...)
}

func Delete(url string, handlers ...Handler) {

	Svr().Delete(url, handlers...)
}

func Patch(url string, handlers ...Handler) {

	Svr().Patch(url, handlers...)
}

func Static(prefix, path string) {

	Svr().Static(prefix, path)
}

func Start() {
	Svr().Start()
}
