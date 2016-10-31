package main

import (
	// "time"
	"github.com/tidwall/buntdb"
	// "github.com/iris-contrib/middleware/basicauth"
	// "github.com/kataras/iris"
)

type Auth struct {
	db *buntdb.DB
}

func (api Auth) install() {
	//authConfig := basicauth.Config{
	//	Users:      map[string]string{"sergeyt": "sergeyt"},
	//	ContextKey: "user_id",            // if you don't set it it's "user"
	//	Expires:    time.Duration(24) * time.Hour,
	//}
	//
	//auth := basicauth.New(authConfig)
	//iris.Party("/api").Use(auth)

	//iris.Post("/api/login", func(ctx *iris.Context) {
	//
	//})
	//iris.Post("/api/logout", func(ctx *iris.Context) {
	//
	//})

}
