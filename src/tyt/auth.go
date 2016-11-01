package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/kataras/iris"
	"github.com/tidwall/buntdb"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	db *buntdb.DB
}

func (api Auth) install() {
	iris.Post("/api/login", func(ctx *iris.Context) {
		authString := ctx.RequestHeader("Authorization")
		auth, err := parseAuthorization(authString)
		if err != nil {
			ctx.EmitError(iris.StatusUnauthorized)
			return
		}
		// TODO bearer is not yet implemented
		if len(auth.token) > 0 {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", authString)

		client := &http.Client{
			Timeout: time.Second * 30,
		}

		res, err := client.Do(req)
		if err != nil {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		data := make(map[string]interface{})
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&data)
		if err != nil {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		login, ok := data["login"].(string)
		if !ok || len(login) == 0 {
			ctx.EmitError(iris.StatusUnauthorized)
			return
		}

		user, err := findUserByLogin(api.db, login)
		if err != nil {
			ctx.EmitError(iris.StatusUnauthorized)
			return
		}

		setUser(ctx, user)
		ctx.JSON(200, "ok")
	})
}

type auth struct {
	login    string
	password string
	token    string
}

var (
	errBadAuthHeader = errors.New("bad auth header")
)

const (
	schemeBasic  = "basic"
	schemeBearer = "bearer"
)

// Retrieves authorization data from given http request.
func parseAuthorization(s string) (*auth, error) {
	if len(s) > 0 {
		var f = strings.Fields(s)
		if len(f) != 2 {
			return nil, errBadAuthHeader
		}
		var scheme = strings.ToLower(f[0])
		var token = f[1]
		switch scheme {
		case schemeBasic:
			var str, err = base64.StdEncoding.DecodeString(token)
			if err != nil {
				return nil, err
			}
			// TODO support realm and auth-params
			var creds = bytes.Split(str, []byte(":"))
			if len(creds) != 2 {
				return nil, errBadAuthHeader
			}
			return &auth{login: string(creds[0]), password: string(creds[1])}, nil
		case schemeBearer:
			return &auth{token: token}, nil
		default:
			return nil, errBadAuthHeader
		}
	}

	return nil, errors.New("no auth token")
}
