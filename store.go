package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/buntdb"
)

func getString(db *buntdb.DB, key string) (string, error) {
	var result string
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		result = val
		return nil
	})
	return result, err
}

func getObject(db *buntdb.DB, key string, out interface{}) error {
	return db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), out)
	})
}

func findUserByLogin(db *buntdb.DB, login string) (*User, error) {
	key := makeLoginKey(login)
	user := &User{}
	err := db.View(func(tx *buntdb.Tx) error {
		id, err := tx.Get(key)
		if err == nil {
			val, err := tx.Get(makeResourceKey("user", id))
			if err != nil {
				return err
			}
			return json.Unmarshal([]byte(val), user)
		}
		// TODO try to do the full scan
		// scan(tx, "user", func(key, val string) bool {
		// 	err := json.Unmarshal([]byte(val), user)
		// 	if err != nil {
		// 		return true
		// 	}
		// 	if strings.ToLower(user.Login) == strings.ToLower(login) {
		// 		return false
		// 	}
		// 	return true
		// })
		// if strings.ToLower(user.Login) == strings.ToLower(login) {
		// 	return nil
		// }
		return err
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

func scan(tx *buntdb.Tx, resource string, handler func(key, val string) bool) error {
	// TODO ascend by prefix
	err := tx.Ascend("", func(key, val string) bool {
		if !strings.HasPrefix(key, resource+"/") {
			return true
		}
		return handler(key, val)
	})
	return err
}

func makeResourceKey(resource, id string) string {
	return fmt.Sprintf("%s/%s", resource, id)
}

func insert(db *buntdb.DB, entity IEntity) error {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	rtype := entity.GetResourceInfo().Type
	key := makeResourceKey(rtype, entity.GetID())
	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(bytes), nil)

		user, ok := entity.(*User)
		if ok {
			onUserInserted(tx, user)
		}

		spec, ok := entity.(*Spectacle)
		if ok {
			onSpectacleInserted(tx, spec)
		}

		return err
	})
}

func onUserInserted(tx *buntdb.Tx, user *User) error {
	// map login to userID
	key := makeLoginKey(user.Login)
	_, _, err := tx.Set(key, user.ID, nil)
	fmt.Printf("map user %s = %s\n", user.Login, user.ID)
	return err
}

func onSpectacleInserted(tx *buntdb.Tx, spec *Spectacle) error {
	key := makeSpectacleKey(spec.Title)
	_, _, err := tx.Set(key, spec.ID, nil)
	return err
}

func makeLoginKey(login string) string {
	return fmt.Sprintf("user:%s", strings.ToLower(login))
}

func makeSpectacleKey(title string) string {
	return fmt.Sprintf("spectacle:%s", title)
}
