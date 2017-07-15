package main

import (
	"fmt"
	"strings"

	"github.com/iris-contrib/plugin/oauth"
	"github.com/kataras/iris"
	"github.com/markbates/goth"
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
		login := strings.ToLower(githubUser.NickName)

		// now find the user in our database
		user, err := findUserByLogin(db, login)
		if user != nil {
			setUser(ctx, user)
			ctx.Redirect("/index.html")
			return
		}

		// try to auto register user if he/she is in list
		userInfo, ok := users[login]
		if !ok {
			ctx.Error(fmt.Sprintf("user not found: %s", err.Error()), iris.StatusUnauthorized)
			return
		}

		user = makeUser(githubUser, userInfo)
		err = insert(db, user)
		if err != nil {
			ctx.Error(err.Error(), iris.StatusInternalServerError)
			return
		}

		setUser(ctx, user)
		ctx.Redirect("/index.html")
	})

	authentication.Fail(func(ctx *iris.Context) {
	})

	iris.Plugins.Add(authentication)
}

func makeUser(gu goth.User, info UserInfo) *User {
	user := &User{
		Name:        info.Name,
		Email:       info.Email,
		Github:      info.Github,
		Role:        info.Role,
		Course:      info.Course,
		Faculty:     info.Faculty,
		Group:       info.Group,
		Description: info.Description,
		Comment:     info.Comment,
		// from github
		Login:     strings.ToLower(gu.NickName),
		AvatarURL: gu.AvatarURL,
		Location:  gu.Location,
	}
	if len(user.Description) == 0 && len(gu.Description) > 0 {
		user.Description = gu.Description
	}
	user.Created("robot")
	return user
}
