package main

import (
	"fmt"
	"github.com/kataras/iris"
)

const (
	keyUserID = "user_id"
)

func getUserID(ctx *iris.Context) string {
	val := ctx.Session().Get(keyUserID)
	if val != nil {
		fmt.Printf("user_id from session: %v", val)
		s, ok := val.(string)
		if ok && len(s) > 0 {
			return s
		}
	}

	val = ctx.Get(keyUserID)
	if val != nil {
		fmt.Printf("user_id from request context: %v", val)
		s, ok := val.(string)
		if ok && len(s) > 0 {
			return s
		}
	}

	s := ctx.GetCookie(keyUserID)
	if len(s) > 0 {
		fmt.Printf("user_id from cookie: %v", val)
		return s
	}

	return "robot"
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
