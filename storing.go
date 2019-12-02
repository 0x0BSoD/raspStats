package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/raspStats/stats"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func DBConn() (*sql.DB, error) {
	createDB := false
	if _, err := os.Stat("./storage"); os.IsNotExist(err) {
		createDB = true
	}

	db, err := sql.Open("sqlite3", "./storage")
	if err != nil {
		return nil, err
	}

	if createDB {
		f, _ := ioutil.ReadFile("./dbInit.sql")
		_, err = db.Exec(string(f))
		if err != nil {
			log.Printf("%q: %s\n", err, f)
			return nil, err
		}

	}

	return db, nil
}

type DBItem struct {
	Uptime  stats.Uptime  `json:"uptime"`
	LoadAvg stats.LoadAvg `json:"load_avg"`
}

func StoreItem(db *sql.DB, data DBItem) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into stats(timestamp, data) values(?, ?)")
	if err != nil {
		return err
	}

	strObj, _ := json.Marshal(data)

	_, err = stmt.Exec(time.Now().Unix(), strObj)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetAllItems(db *sql.DB) (string, error) {
	rows, err := db.Query("select timestamp, data from stats")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var timestamp time.Time
		var data string
		err = rows.Scan(&timestamp, &data)
		if err != nil {
			return "", err
		}
		fmt.Println(timestamp, data)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	return "", nil
}
