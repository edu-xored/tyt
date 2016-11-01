package main

import (
	"fmt"
	"github.com/iris-contrib/plugin/oauth"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
)

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

		setUser(ctx, user)
		ctx.Redirect("/index.html")
	})

	authentication.Fail(func(ctx *iris.Context) {
	})

	iris.Plugins.Add(authentication)
}
