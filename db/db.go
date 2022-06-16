package db

import (
	"database/sql"
	// "fmt"
	"log"
	// "os"
)

var DB *sql.DB

func Init() {
	// err := os.WriteFile("filename.txt", []byte("Hello"), 0755)
	// if err != nil {
	// 	fmt.Printf("Unable to write file: %v", err)
	// }
	database, connectionErr := sql.Open("sqlite3", "database.db")
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}
	DB = database

}
