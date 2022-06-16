package main

import (
<<<<<<< HEAD
	// "fmt"
	// "log"
	// "strconv"
=======
	"fmt"
	"log"
	"strconv"
>>>>>>> 962bbc009fcdc18a00d8265444f3ec19791a1fce

	// "text/template"

	"net/http"

	"forumynov.com/db"
	"forumynov.com/routes"
	"github.com/gorilla/context"
	_ "github.com/mattn/go-sqlite3"
)

var err error

func main() {
<<<<<<< HEAD

	db.Init()
=======
	db.Init()
	rows, err := db.DB.Query("SELECT (id, username, password) FROM users WHERE username=\"?\"")
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
>>>>>>> 962bbc009fcdc18a00d8265444f3ec19791a1fce

	http.HandleFunc("/signup", routes.SignupPage)
	http.HandleFunc("/login", routes.LoginPage)
	http.HandleFunc("/logout", routes.LogoutPage)
	http.HandleFunc("/", routes.HomePage)
	http.HandleFunc("/logged", routes.HomePageLogged)

	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	defer db.DB.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
