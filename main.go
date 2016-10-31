package main

import (
	"log"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/plugin/oauth"
	"fmt"
)

func main() {
	viper.SetConfigName("tyt")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// FIXME
	iris.Config.VScheme = "http://"
	iris.Config.VHost = "localhost:8080"

	initOAuth()
	installAPI(db)
	staticContent()

	iris.Listen(":8080")
}

func installAPI(db *buntdb.DB) {
	iris.Use(logger.New())
	iris.Use(cors.Default())

	Auth { db: db }.install()

	// users API
	// TODO validation, new user must have unique email and login
	API {
		db: db,
		resource:"user",
		collection:"users",
		factory: func() IEntity { return &User{} },
	}.install()

	// teams API
	API {
		db: db,
		resource:"team",
		collection:"teams",
		factory: func() IEntity { return &Team{} },
	}.install()

	// events API
	API {
		db: db,
		resource:"event",
		collection:"events",
		factory: func() IEntity { return &Event{} },
	}.install()

	// TODO api to change user password
}

func staticContent() {
	iris.Static("/public", "./public/", 1)

	iris.Get("/", func(ctx *iris.Context) {
		ctx.ServeFile("./index.html", false)
	});

	iris.Get("/index.html", func(ctx *iris.Context) {
		ctx.ServeFile("./index.html", false)
	});

	iris.Get("/login.html", func(ctx *iris.Context) {
		ctx.ServeFile("./login.html", false)
	});
	iris.Get("/statistics", func(ctx *iris.Context) {
		ctx.ServeFile("./statistics.html", false)
	});
}

func initOAuth() {
	clientID := viper.Get("github.client_id")
	secret := viper.Get("github.client_secret")

	// register your auth via configs, providers with non-empty
	// values will be registered to goth automatically by Iris
	oauthConfig := oauth.Config{
		Path: "/oauth",
		GithubKey:   clientID.(string),
		GithubSecret: secret.(string),
	}

	authentication := oauth.New(oauthConfig)
	iris.Plugins.Add(authentication)

	// came from host:port/oauth/:provider
	// this is the handler inside host:port/oauth/:provider/callback
	// you can do redirect to the authenticated url or whatever you want to do
	authentication.Success(func(ctx *iris.Context) {
		user := authentication.User(ctx) // returns the goth.User
		fmt.Printf("github user: %v", user)
	})
	authentication.Fail(func(ctx *iris.Context){
	})
}
