package main

import (
	"fmt"
	"log"

	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
	"strings"
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

	getCurrentUser = makeGetCurrentUser(db)

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
