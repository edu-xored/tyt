package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/plugin/oauth"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
	"strings"
)

const (
	keyUserID = "user_id"
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

	initOAuth(db)
	installAPI(db)
	staticContent()

	iris.Listen(":8080")
}

func installAPI(db *buntdb.DB) {
	iris.Use(logger.New())
	iris.Use(cors.Default())

	Auth{db: db}.install()

	iris.Get("/api/me", func(ctx *iris.Context) {
		userID := getUserID(ctx)
		if len(userID) == 0 {
			ctx.NotFound()
			return
		}
		fmt.Printf("user_id = %s\n", userID)
		user, err := findUserByID(db, userID)
		if err != nil {
			ctx.NotFound()
			return
		}
		ctx.JSON(200, user)
	})

	// users API
	// TODO validation, new user must have unique email and login
	API{
		db:         db,
		resource:   "user",
		collection: "users",
		factory:    func() IEntity { return &User{} },
		onCreate: func(tx *buntdb.Tx, r IEntity) error {
			user := r.(*User)
			// map login to userID
			key := fmt.Sprintf("user:%s", strings.ToLower(user.Login))
			_, _, err := tx.Set(key, user.ID, nil)
			fmt.Printf("map user %s = %s", user.Login, user.ID)
			return err
		},
	}.install()

	// teams API
	API{
		db:         db,
		resource:   "team",
		collection: "teams",
		factory:    func() IEntity { return &Team{} },
	}.install()

	// events API
	API{
		db:         db,
		resource:   "event",
		collection: "events",
		factory:    func() IEntity { return &Event{} },
	}.install()

	// TODO api to change user password
}

func staticContent() {
	iris.Static("/public", "./public/", 1)

	index := func(ctx *iris.Context) {
		// TODO better approach to authorize
		userID := getUserID(ctx)
		if len(userID) == 0 {
			ctx.Redirect("/login.html")
			return
		}
		ctx.ServeFile("./index.html", false)
	}

	iris.Get("/login.html", func(ctx *iris.Context) {
		ctx.ServeFile("./login.html", false)
	})

	iris.Get("/", index)
	iris.Get("/index.html", index)

	iris.Get("/statistics", func(ctx *iris.Context) {
		// TODO better approach to authorize
		user_id := ctx.GetCookie(keyUserID)
		if len(user_id) == 0 {
			ctx.Redirect("/login.html")
			return
		}
		ctx.ServeFile("./statistics.html", false)
	})
}

func initOAuth(db *buntdb.DB) {
	clientID := viper.Get("github.client_id")
	secret := viper.Get("github.client_secret")

	// register your auth via configs, providers with non-empty
	// values will be registered to goth automatically by Iris
	oauthConfig := oauth.Config{
		Path:         "/oauth",
		GithubKey:    clientID.(string),
		GithubSecret: secret.(string),
	}

	authentication := oauth.New(oauthConfig)

	// came from host:port/oauth/:provider
	// this is the handler inside host:port/oauth/:provider/callback
	// you can do redirect to the authenticated url or whatever you want to do
	authentication.Success(func(ctx *iris.Context) {
		githubUser := authentication.User(ctx) // returns the goth.User
		//json, err := json.MarshalIndent(githubUser, "", "  ")
		//if err == nil {
		//	fmt.Printf("github user: %s\n", string(json))
		//}

		// now find the user in our database
		user, err := findUserByLogin(db, githubUser.NickName)
		if user == nil {
			ctx.Error(fmt.Sprintf("user not found: %s", err.Error()), iris.StatusUnauthorized)
			return
		}

		// TODO generate token instead of user id
		ctx.Set(keyUserID, user.ID)
		ctx.SetCookieKV(keyUserID, user.ID)
		ctx.Redirect("/index.html")
	})

	authentication.Fail(func(ctx *iris.Context) {
	})

	iris.Plugins.Add(authentication)
}

func findUserByLogin(db *buntdb.DB, login string) (*User, error) {
	key := fmt.Sprintf("user:%s", strings.ToLower(login))
	user := &User{}
	err := db.View(func(tx *buntdb.Tx) error {
		id, err := tx.Get(key)
		if err != nil {
			return err
		}
		val, err := tx.Get(fmt.Sprintf("user/%s", id))
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), user)
	})
	if err != nil {
		fmt.Printf("find user failed: %v", err)
		return nil, err
	}
	return user, nil
}

func findUserByID(db *buntdb.DB, id string) (*User, error) {
	user := &User{}
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(fmt.Sprintf("user/%s", id))
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), user)
	})
	if err != nil {
		fmt.Printf("find user failed: %v", err)
		return nil, err
	}
	return user, nil
}
