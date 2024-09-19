package db

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/structs"
	"github.com/gertd/go-pluralize"
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

			var err error
			ret, err = creator(targetName)

			if err != nil {
				log.Fatal("Unable to create DB connection ", err)
				return nil
			}

			listConn[targetName] = ret
		}
	}

	return ret
}

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var pluralClient = pluralize.NewClient()

func ToSnakeCase(str string) string {

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetIDValue(obj interface{}) string {

	list := structs.Map(obj)

	id, exist := list["ID"]

	if exist {

		return id.(string)
	}

	return ""
}

func GetTableName(obj interface{}) string {

	ret := ""
	targetObj := obj

	if reflect.TypeOf(obj).Kind() == reflect.Pointer {

		if reflect.TypeOf(obj).Elem().Kind() == reflect.Array ||
			reflect.TypeOf(obj).Elem().Kind() == reflect.Slice ||
			reflect.TypeOf(obj).Elem().Kind() == reflect.Map {

			targetObj = reflect.New(reflect.TypeOf(obj).Elem().Elem()).Interface()
		}
	} else if reflect.TypeOf(obj).Kind() == reflect.Array ||
		reflect.TypeOf(obj).Kind() == reflect.Slice ||
		reflect.TypeOf(obj).Kind() == reflect.Map {

		targetObj = reflect.New(reflect.TypeOf(obj).Elem()).Interface()
	}

	informer, ok := targetObj.(TableNameInformer)

	if ok {
		ret = informer.TableName()
	}

	if ret == "" {

		ret = ToSnakeCase(pluralClient.Plural(structs.Name(targetObj)))
	}

	return ret
}

func Asc(fieldName string) FieldOrder {

	return FieldOrder{
		Field: fieldName,
		Order: "ASC",
	}
}

func Desc(fieldName string) FieldOrder {

	return FieldOrder{
		Field: fieldName,
		Order: "DESC",
	}
}

func Save(obj interface{}) error {

	return Conn().Save(obj)
}

func UpdateMany(table interface{}, filter interface{}, update interface{}) error {

	return Conn().Update(table, filter, update)
}

func Delete(obj interface{}) error {

	return Conn().Delete(obj)
}

func DeleteMany(table interface{}, filter interface{}) error {

	return Conn().DeleteMany(table, filter)
}

func FindById(out interface{}, id string) error {
	return Conn().FindById(out, id)
}

func FindOne(out interface{}, filter interface{}, options ...FindOption) error {

	return Conn().FindOne(out, filter, options...)
}

func Find(out interface{}, filter interface{}, options ...FindOption) error {

	return Conn().Find(out, filter, options...)
}
