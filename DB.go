package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func openDatabase() error {
	var err error
	DB, err = sql.Open("postgres", "user=postgres password=postgressUserPa55 dbname=Aces-Whist sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}

func closeDatabase() error {
	return DB.Close()
}