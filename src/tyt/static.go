package main

import (
	"github.com/kataras/iris"
)

func staticContent() {
	// public content
	iris.Static("/public", "./public/", 1)

	iris.Get("/login.html", func(ctx *iris.Context) {
		ctx.ServeFile("./login.html", false)
	})

	iris.Get("/docs", func(ctx *iris.Context) {
		ctx.ServeFile("./docs/index.html", false)
	})
	iris.Get("/docs/index.html", func(ctx *iris.Context) {
		ctx.ServeFile("./docs/index.html", false)
	})

	protected := func(handler iris.HandlerFunc) iris.HandlerFunc {
		return func(ctx *iris.Context) {
			// TODO better approach to authorize
			user := getCurrentUser(ctx)
			if user == nil {
				ctx.Redirect("/login.html")
				return
			}
			handler(ctx)
		}
	}

	index := protected(func(ctx *iris.Context) {
		ctx.ServeFile("./index.html", false)
	})

	iris.Get("/", index)
	iris.Get("/index.html", index)

	iris.Get("/statistics.html", protected(func(ctx *iris.Context) {
		ctx.ServeFile("./statistics.html", false)
	}))
}
