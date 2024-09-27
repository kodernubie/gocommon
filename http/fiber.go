package http

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kodernubie/gocommon/conf"
	"github.com/kodernubie/gocommon/log"
)

//---------------------------------------------------------------------

type FiberContext struct {
	ctx *fiber.Ctx
}

func (o *FiberContext) Raw() interface{} {

	return o.ctx
}

func (o *FiberContext) Bind(target interface{}) error {

	return json.Unmarshal(o.ctx.Body(), target)
}

func (o *FiberContext) Param(key string, defaultVal ...string) string {

	return o.ctx.Params(key, defaultVal...)
}

func (o *FiberContext) Query(key string, defaultVal ...string) string {

	return o.ctx.Query(key, defaultVal...)
}

func (o *FiberContext) QueryInt(key string, defaultVal ...int) int {

	ret := 0
	retStr := o.ctx.Query(key, "")

	if retStr == "" {

		if len(defaultVal) > 0 {
			ret = defaultVal[0]
		}
	} else {

		ret, _ = strconv.Atoi(retStr)
	}

	return ret
}

func (o *FiberContext) Body() []byte {

	return o.ctx.Body()
}

func (o *FiberContext) FormFile(key string) (*multipart.FileHeader, error) {

	return o.ctx.FormFile(key)
}

func (o *FiberContext) SaveFile(key string, path string) error {

	fl, err := o.ctx.FormFile(key)

	if err != nil {
		return err
	}

	return o.ctx.SaveFile(fl, path)
}

func (o *FiberContext) FormValue(key string, defaultValue ...string) string {

	return o.ctx.FormValue(key, defaultValue...)
}

func (o *FiberContext) Status(code int) {

	o.ctx.Status(code)
}

func (o *FiberContext) Reply(data interface{}) error {

	return o.ctx.JSON(ReplyObj{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func (o *FiberContext) ReplyPage(data interface{}, pageNo int, hasNext bool) error {

	return o.ctx.JSON(ReplyPageObj{
		ReplyObj: ReplyObj{
			Code: 0,
			Msg:  "success",
			Data: data,
		},
		PageNo:  pageNo,
		HasNext: hasNext,
	})
}

func (o *FiberContext) ReplyRaw(httpCode int, data ...[]byte) error {

	if len(data) > 0 {
		return o.ctx.Status(httpCode).Send(data[0])
	}

	return o.ctx.SendStatus(httpCode)
}

func (o *FiberContext) Error(code int, msg string, data ...interface{}) error {

	retObj := ReplyObj{
		Code: code,
		Msg:  msg,
	}

	if len(data) > 0 {
		retObj.Data = data[0]
	}

	return o.ctx.Status(400).JSON(retObj)
}

func (o *FiberContext) SendFile(fileName string) error {

	return o.ctx.SendFile(fileName)
}

func (o *FiberContext) Next() error {

	return o.ctx.Next()
}

func (o *FiberContext) ReqHeader(key string) string {

	list, exist := o.ctx.GetReqHeaders()[key]

	if !exist || len(list) == 0 {
		return ""
	}

	return list[0]
}

func (o *FiberContext) SetHeader(key, value string) {

	o.ctx.Set(key, value)
}

func (o *FiberContext) GetHeader(key string) string {

	return o.ctx.Get(key)
}

func (o *FiberContext) RequestURL() string {

	return string(o.ctx.OriginalURL())
}

func (o *FiberContext) Redirect(to string) error {

	return o.ctx.Redirect(to, http.StatusMovedPermanently)
}

func (o *FiberContext) Download(fileName string) error {

	return o.ctx.Download(fileName)
}

func (o *FiberContext) Set(key string, val interface{}) {

	o.ctx.Context().SetUserValue(key, val)
}

func (o *FiberContext) Get(key string) interface{} {

	return o.ctx.Context().UserValue(key)
}

//---------------------------------------------------------------------------------------------

type FiberServer struct {
	app        *fiber.App
	configName string
}

func (o *FiberServer) convertHandler(handlers ...Handler) []fiber.Handler {

	ret := []fiber.Handler{}

	for _, handler := range handlers {

		ret = append(ret, toFiberHandler(handler))
	}

	return ret
}

func toFiberHandler(handler Handler) fiber.Handler {

	return func(c *fiber.Ctx) error {

		defer func() {

			err := recover()

			if err != nil {
				log.Errorf("recover from http handler error :", err)
			}
		}()

		return handler(&FiberContext{
			ctx: c,
		})
	}
}

func (o *FiberServer) Raw() interface{} {

	return o.app
}

func (o *FiberServer) Get(path string, handlers ...Handler) {

	o.app.Get(path, o.convertHandler(handlers...)...)
}

func (o *FiberServer) Post(path string, handlers ...Handler) {

	o.app.Post(path, o.convertHandler(handlers...)...)
}

func (o *FiberServer) Put(path string, handlers ...Handler) {

	o.app.Put(path, o.convertHandler(handlers...)...)
}

func (o *FiberServer) Patch(path string, handlers ...Handler) {

	o.app.Patch(path, o.convertHandler(handlers...)...)
}

func (o *FiberServer) Delete(path string, handlers ...Handler) {

	o.app.Delete(path, o.convertHandler(handlers...)...)
}

func (o *FiberServer) Static(prefix string, path string) {

	o.app.Static(prefix, path)
}

func (o *FiberServer) Start() {

	port := conf.Str("HTTP_" + o.configName + "_PORT")

	if port == "" {
		port = "8080"
	}
	addr := conf.Str("HTTP_" + o.configName + "_ADDRES")

	o.app.Listen(addr + ":" + port)
}

func init() {

	RegCreator("fiber", func(name string) (Server, error) {

		ret := &FiberServer{}

		ret.app = fiber.New()

		ret.configName = name

		// cors
		hasCorsConfig := false
		corsConfig := cors.Config{}

		if conf.Str("HTTP_"+name+"_CORS_ORIGIN") != "" {
			corsConfig.AllowOrigins = conf.Str("HTTP_" + name + "_CORS_ORIGIN")
			hasCorsConfig = true
		}

		if conf.Str("HTTP_"+name+"_CORS_HEADER") != "" {
			corsConfig.AllowHeaders = conf.Str("HTTP_" + name + "_CORS_HEADER")
			hasCorsConfig = true
		}

		if hasCorsConfig {
			ret.app.Use(cors.New(corsConfig))
		} else {
			ret.app.Use(cors.New())
		}

		return ret, nil
	})
}
