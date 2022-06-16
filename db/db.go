package db

import (
	"database/sql"
	"log"
)

var err error

var DB *sql.DB

func Init() {

	database, connectionErr := sql.Open("sqlite3", "database.db")
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}
	DB = database

}
