package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/tidwall/buntdb"
	"strings"
)

const contentJSON = "application/json"

type API struct {
	db         *buntdb.DB
	resource   string // singular name
	collection string
	factory    func() IEntity
	onCreate   func(tx *buntdb.Tx, resource IEntity) error
}

func (api API) install() {
	collectionPath := fmt.Sprintf("/api/%s", api.collection)
	resourcePath := fmt.Sprintf("/api/%s/:id", api.resource)
	iris.Post(collectionPath, api.makeCreateHandler())
	iris.Get(collectionPath, api.makeListHandler())
	iris.Get(resourcePath, api.makeGetHandler())
	iris.Put(resourcePath, api.makePutHandler())
	iris.Delete(resourcePath, api.makeDeleteHandler())
}

func (api API) decode(val string) (interface{}, error) {
	resource := api.factory()
	err := json.Unmarshal([]byte(val), resource)
	return resource, err
}

func (api API) makeListHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		err := api.db.View(func(tx *buntdb.Tx) error {
			// TODO use string builder
			comma := false
			count := 0
			str := "["
			err := tx.Ascend("", func(key, val string) bool {
				if !strings.HasPrefix(key, api.resource+"/") {
					return true
				}
				if comma {
					str += ","
				}
				str += val
				comma = true
				count = count + 1
				return true
			})
			if err != nil {
				return err
			}
			str += "]"

			list := make([]map[string]interface{}, count)
			err = json.Unmarshal([]byte(str), &list)
			if err != nil {
				return err
			}

			ctx.JSON(200, list)

			// ctx.RenderWithStatus(200, contentJSON, result)

			return nil
		})
		if err != nil {
			sendError(ctx, err)
		}
	}
}

func (api API) makeResourceKey(ctx *iris.Context) string {
	id := ctx.Param("id")
	return fmt.Sprintf("%s/%s", api.resource, id)
}

func (api API) makeGetHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		key := api.makeResourceKey(ctx)
		err := api.db.View(func(tx *buntdb.Tx) error {
			val, err := tx.Get(key)
			if err != nil {
				return err
			}
			// TODO send JSON as is without deserialization
			// ctx.RenderWithStatus(200, contentJSON, val)

			res, err := api.decode(val)
			if err != nil {
				return err
			}
			ctx.JSON(200, res)

			return nil
		})
		if err != nil {
			sendError(ctx, err)
		}
	}
}

func getUserID(ctx *iris.Context) string {
	val := ctx.Get(keyUserID)
	if val != nil {
		return val.(string)
	}
	s := ctx.GetCookie(keyUserID)
	if len(s) > 0 {
		return s
	}
	return "robot"
}

func (api API) readResource(ctx *iris.Context) (IEntity, error) {
	resource := api.factory()

	if err := ctx.ReadJSON(resource); err != nil {
		logError(ctx, "can not parse")
		return nil, err
	}

	return resource, nil
}

func logError(ctx *iris.Context, message string) {
	payload := string(ctx.Request.Body())
	fmt.Printf("%s: %s", message, payload)
}

func (api API) makeCreateHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		err := api.db.Update(func(tx *buntdb.Tx) error {

			resource, err := api.readResource(ctx)
			if err != nil {
				return err
			}

			resource.Created(getUserID(ctx))

			bytes, err := json.Marshal(resource)
			if err != nil {
				logError(ctx, "can not marshal")
				return err
			}

			key := fmt.Sprintf("%s/%s", api.resource, resource.GetID())
			_, _, err = tx.Set(key, string(bytes), nil)

			if api.onCreate != nil {
				err = api.onCreate(tx, resource)
				if err != nil {
					return err
				}
			}

			if err != nil {
				return err
			}

			ctx.JSON(200, resource)

			return nil
		})
		if err != nil {
			sendError(ctx, err)
		}
	}
}

func (api API) makePutHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		id := ctx.Param("id")
		key := api.makeResourceKey(ctx)

		err := api.db.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Get(key)
			if err != nil {
				return err
			}

			resource, err := api.readResource(ctx)
			if err != nil {
				return err
			}

			if resource.GetID() != id {
				return fmt.Errorf("payload must have the same id as in URL, but was %s", resource.GetID())
			}

			resource.Updated(getUserID(ctx))

			bytes, err := json.Marshal(resource)
			if err != nil {
				logError(ctx, "can not marshal")
				return err
			}

			_, _, err = tx.Set(key, string(bytes), nil)

			if err != nil {
				return err
			}

			ctx.JSON(200, resource)

			return nil
		})

		if err != nil {
			sendError(ctx, err)
		}
	}
}

func (api API) makeDeleteHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		key := api.makeResourceKey(ctx)
		err := api.db.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Delete(key)
			return err
		})
		if err != nil {
			sendError(ctx, err)
		} else {
			// TODO send just status
			ctx.JSON(200, "ok")
		}
	}
}

func sendError(ctx *iris.Context, err error) {
	fmt.Printf("error: %v", err)
	// TODO classify errors
	ctx.Error(err.Error(), 404)
}
