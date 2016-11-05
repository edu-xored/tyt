package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/buntdb"
)

func findUserByLogin(db *buntdb.DB, login string) (*User, error) {
	key := fmt.Sprintf("user:%s", strings.ToLower(login))
	user := &User{}
	err := db.View(func(tx *buntdb.Tx) error {
		id, err := tx.Get(key)
		if err == nil {
			val, err := tx.Get(fmt.Sprintf("user/%s", id))
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

		return err
	})
}

func onUserInserted(tx *buntdb.Tx, user *User) error {
	// map login to userID
	key := fmt.Sprintf("user:%s", strings.ToLower(user.Login))
	_, _, err := tx.Set(key, user.ID, nil)
	fmt.Printf("map user %s = %s\n", user.Login, user.ID)
	return err
}
