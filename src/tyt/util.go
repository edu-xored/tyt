package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/tidwall/buntdb"
)

const (
	keyUserID = "user_id"
)

var getCurrentUser func(ctx *iris.Context) *User

func makeGetCurrentUser(db *buntdb.DB) func(ctx *iris.Context) *User {
	return func(ctx *iris.Context) *User {
		id := getUserID(ctx)
		user, err := findUserByID(db, id)
		if err != nil {
			return nil
		}
		return user
	}
}

func getUserID(ctx *iris.Context) string {
	val := ctx.Session().Get(keyUserID)
	if val != nil {
		fmt.Printf("user_id from session: %v\n", val)
		s, ok := val.(string)
		if ok && len(s) > 0 {
			return s
		}
	}

	val = ctx.Get(keyUserID)
	if val != nil {
		fmt.Printf("user_id from request context: %v\n", val)
		s, ok := val.(string)
		if ok && len(s) > 0 {
			return s
		}
	}

	s := ctx.GetCookie(keyUserID)
	if len(s) > 0 {
		fmt.Printf("user_id from cookie: %v\n", val)
		return s
	}

	return "robot"
}

func setUser(ctx *iris.Context, user *User) {
	// TODO generate token instead of user id
	ctx.Set(keyUserID, user.ID)
	ctx.Session().Set(keyUserID, user.ID)
	ctx.SetCookieKV(keyUserID, user.ID)
}

func logError(ctx *iris.Context, message string) {
	payload := string(ctx.Request.Body())
	fmt.Printf("%s: %s", message, payload)
}

func sendError(ctx *iris.Context, err error) {
	fmt.Printf("error: %v", err)
	// TODO classify errors
	ctx.Error(err.Error(), 404)
}
