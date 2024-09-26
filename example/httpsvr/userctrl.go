package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/kodernubie/gocommon/http"
	"github.com/kodernubie/gocommon/log"
	"github.com/oklog/ulid/v2"
)

var userCache map[string]*User = map[string]*User{}
var tokenMap map[string]string = map[string]string{} // example of token repo, use jwt in real app

func init() {

	http.Post("/user/login", login)

	http.Get("/user/me", checkToken, me) // example of multi handler for 1 api, will be called sequentially

	http.Post("/user", createUser)
	http.Get("/user", searchUser) // example api using query

	http.Post("/user/document", upload)  // example api to upload file
	http.Get("/user/document", download) // example api to download file

	http.Get("/user/:id", getUser) // example api with path parameter
	http.Put("/user/:id", updateUser)
	http.Delete("/user/:id", delUser)

}

func getUserByName(name string) *User {

	var ret *User

	for _, user := range userCache {

		if user.Name == name {
			ret = user
			break
		}
	}

	return ret
}

// example of POST method that accept json payload
func createUser(ctx http.Context) error {

	req := &UserReq{}
	err := ctx.Bind(req)

	if err != nil {
		return ctx.Error(1001, "Unable to parse request "+err.Error())
	}

	existing := getUserByName(req.Name)

	if existing != nil {
		return ctx.Error(1001, "User already exist")
	}

	user := &User{}
	copier.Copy(&user, req)

	user.ID = ulid.Make().String()
	userCache[user.ID] = user

	ret := UserRes{}
	copier.Copy(&ret, user)

	return ctx.Reply(ret)
}

// example of GET method with pat parameter
func getUser(ctx http.Context) error {

	user, exist := userCache[ctx.Param("id")]

	if !exist {

		return ctx.Error(1001, "user not exist")
	}

	ret := UserRes{}
	copier.Copy(&ret, user)

	return ctx.Reply(ret)
}

// example of PUT method with pat parameter and accept JSON payload
func updateUser(ctx http.Context) error {

	req := &UserReq{}
	err := ctx.Bind(req)

	if err != nil {
		return ctx.Error(1001, "Unable to parse request "+err.Error())
	}

	user, exist := userCache[ctx.Param("id")]

	if !exist {
		return ctx.Error(1001, "user not exist")
	}

	user.Name = req.Name
	user.Password = req.Password

	ret := UserRes{}
	copier.Copy(&ret, user)

	return ctx.Reply(ret)
}

// example DELETE method with path parameter
func delUser(ctx http.Context) error {

	delete(userCache, ctx.Param("id"))
	return ctx.Reply(true)
}

func searchUser(ctx http.Context) error {

	filter := ctx.Query("filter")
	ret := []UserRes{}

	for _, user := range userCache {

		if strings.Contains(user.Name, filter) {

			obj := UserRes{}
			copier.Copy(&obj, user)

			ret = append(ret, obj)
		}
	}

	return ctx.ReplyPage(ret, 1, 1)
}

// example login that respin with access token than will be used in restricted api
func login(ctx http.Context) error {

	req := &UserReq{}
	err := ctx.Bind(req)

	if err != nil {
		return ctx.Error(1001, "Unable to parse request "+err.Error())
	}

	user := getUserByName(req.Name)

	if user == nil || req.Password != user.Password {
		return ctx.Error(1001, "User not found or password mismatch")
	}

	ret := UserLoginRes{}
	copier.Copy(&ret, user)

	// just example, dont use this in real app
	token := sha256.New()
	token.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	ret.AccessToken = fmt.Sprintf("%x", token.Sum(nil))

	tokenMap[ret.AccessToken] = user.ID

	return ctx.Reply(ret)
}

// example of checking authorization from request header
// and saving custom data in context that can be used in another handler
func checkToken(ctx http.Context) error {

	token := ctx.ReqHeader("Authorization")

	if token == "" {
		return ctx.Error(1002, "not authorized")
	}

	userId, exist := tokenMap[token[7:]]

	if !exist {
		return ctx.Error(1002, "not authorized")
	}

	ctx.Set("userId", userId)

	return ctx.Next()
}

// example of handler that use custom data that already set by previous handler
func me(ctx http.Context) error {

	user, exist := userCache[fmt.Sprintf("%v", ctx.Get("userId"))]

	if !exist {
		return ctx.Error(1002, "user not found")
	}

	return ctx.Reply(user)
}

func upload(ctx http.Context) error {

	uploadDoc, err := ctx.FormFile("doc")

	if err != nil {
		log.Info("error :", err)
		ctx.Error(1013, "Error :", err)
	}

	ctx.SaveFile("doc", "./uploaded.pdf")

	return ctx.Reply(uploadDoc.Filename)
}

func download(ctx http.Context) error {

	return ctx.Download("./sample.pdf")
}
