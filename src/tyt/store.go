package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"strings"
)

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
