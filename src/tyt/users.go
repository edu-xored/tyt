package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type UserInfo struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Github      string `json:"github,omitempty"`
	Role        string `json:"role,omitempty"`
	Course      int32  `json:"course"`
	Faculty     string `json:"faculty,omitempty"`
	Group       int32  `json:"group,omitempty"`
	Description string `json:"description"`
	Comment     string `json:"comment"`
}

var users = make(map[string]UserInfo)

func loadUsers() {
	data, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var arr []UserInfo
	err = json.Unmarshal(data, &arr)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, u := range arr {
		i := strings.LastIndex(u.Github, "/")
		login := strings.ToLower(u.Github[i+1:])
		users[strings.ToLower(u.Github)] = u
		users[login] = u
	}
}
