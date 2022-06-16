package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	// "text/template"

	"net/http"

	"forumynov.com/routes"
	"github.com/gorilla/context"
	_ "github.com/mattn/go-sqlite3"
)

var err error

var db *sql.DB

func main() {

	database, connectionErr := sql.Open("sqlite3", "./database.db")
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}
	defer database.Close()
	db = database

	rows, err := database.Query("SELECT id, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rows)

	var id int
	var username string
	var password string

	for rows.Next() {
		rows.Scan(&id, &username, &password)
		fmt.Println(strconv.Itoa(id) + ": " + username + " " + password)
	}

	http.HandleFunc("/signup", routes.SignupPage)
	http.HandleFunc("/login", routes.LoginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func deleteUser(id int) {
	// regarder si l'id correspond a l'utilisateur connecté

	_, err := db.Exec("DELTE FROM users WHERE id=\"?\";", id)
	checkErr(err)

}
