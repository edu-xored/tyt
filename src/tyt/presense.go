package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kataras/iris"
	"github.com/tidwall/buntdb"
)

func initPresenceAPI(db *buntdb.DB) {
	iris.Post("/api/iamhere", func(ctx *iris.Context) {
		user := getCurrentUser(ctx)
		if user == nil {
			ctx.EmitError(iris.StatusUnauthorized)
			return
		}

		// TODO block by X-Real-IP

		input := &struct {
			SpectacleID string `json:"spectacle_id"`
		}{}
		if err := ctx.ReadJSON(input); err != nil {
			ctx.EmitError(iris.StatusBadRequest)
			return
		}

		fmt.Printf("presence from %s for lecture %s\n", user.Login, input.SpectacleID)

		spec := &Spectacle{}
		err := getObject(db, makeResourceKey("spectacle", input.SpectacleID), spec)
		if err != nil || spec == nil {
			fmt.Printf("lecture %s not found\n", input.SpectacleID)
			ctx.EmitError(iris.StatusBadRequest)
			return
		}

		// check report is in time
		now := time.Now().UTC()
		end := spec.Start.Add(time.Duration(spec.Duration) * time.Hour)
		if !within(spec.Start, end, now) {
			ctx.EmitError(iris.StatusBadRequest)
			return
		}

		event := findPresenceEvent(db, spec.ID)
		if event != nil {
			fmt.Printf("presence to lecture %s already reported by %s\n", input.SpectacleID, user.Login)
			ctx.EmitError(iris.StatusBadRequest)
			return
		}

		event = &Event{
			UserID:      user.ID,
			Type:        EventPresence,
			Message:     "I am here!",
			Start:       now,
			Duration:    spec.Duration,
			SpectacleID: spec.ID,
		}
		event.Created(user.ID)
		err = insert(db, event)
		if err != nil {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		ctx.JSON(200, "ok")
	})
}

func findPresenceEvent(db *buntdb.DB, specID string) *Event {
	var result *Event
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		scan(tx, "event", func(key, val string) bool {
			evt := &Event{}
			err = json.Unmarshal([]byte(val), evt)
			if err != nil {
				return false
			}
			if evt.Type == EventPresence && evt.SpectacleID == specID {
				result = evt
				return false
			}
			return true
		})
		return err
	})
	if err != nil {
		return nil
	}
	return result
}

func within(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
