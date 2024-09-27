package http

import "mime/multipart"

type Server interface {
	Get(path string, handler ...Handler)
	Post(path string, handler ...Handler)
	Put(path string, handler ...Handler)
	Delete(path string, handler ...Handler)
	Patch(path string, handler ...Handler)
	Static(prefix string, path string)
	Raw() interface{}
	Start()
}

type Creator func(configName string) (Server, error)

type Context interface {
	Bind(target interface{}) error
	Param(key string, defaultVal ...string) string
	Query(key string, defaultValue ...string) string
	QueryInt(key string, defaultValue ...int) int
	Body() []byte
	FormFile(key string) (*multipart.FileHeader, error)
	SaveFile(key string, path string) error
	FormValue(key string, defaultValue ...string) string
	Status(code int)
	Next() error
	ReqHeader(key string) string
	SetHeader(key, value string)
	GetHeader(key string) string
	RequestURL() string
	Redirect(to string) error
	Reply(data interface{}) error
	ReplyPage(data interface{}, pageNo int, hasNext bool) error
	ReplyRaw(httpCode int, data ...[]byte) error
	Error(code int, msg string, data ...interface{}) error
	Download(fileName string) error
	Set(key string, val interface{})
	Get(key string) interface{}
	Raw() interface{}
}

type Handler func(ctx Context) error

type ReplyObj struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type ReplyPageObj struct {
	ReplyObj
	PageNo  int  `json:"pageNo"`
	HasNext bool `json:"hasNext"`
}
