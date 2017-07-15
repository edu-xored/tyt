package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"fmt"
	"time"

	"github.com/tidwall/buntdb"
)

type lecture struct {
	Title     string `json:"title"`
	Start     string `json:"start"`
	Presenter string `json:"presenter"`
}

func loadLectures(db *buntdb.DB) {
	content, err := ioutil.ReadFile("./data/lectures.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var list []lecture
	err = json.Unmarshal(content, &list)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, it := range list {
		err = initLecture(db, it)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func initLecture(db *buntdb.DB, lec lecture) error {
	start, err := time.ParseInLocation("2006-01-02 15:04", lec.Start, time.Local)
	if err != nil {
		return err
	}

	key := makeSpectacleKey(lec.Title)
	_, err = getString(db, key)
	if err == nil {
		// TODO update existing lecture
		return nil
	}

	spec := &Spectacle{
		Type:          "lecture",
		Title:         lec.Title,
		PresenterName: lec.Presenter,
		Start:         start.UTC(),
		Duration:      1,
	}
	spec.Created("robot")
	err = insert(db, spec)

	if err == nil {
		fmt.Printf("lecture '%s' initialized\n", lec.Title)
	}

	return err
}
