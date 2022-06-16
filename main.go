package main

import (
	// "fmt"
	// "log"
	// "strconv"

	// "text/template"

	"net/http"

	"forumynov.com/db"
	"forumynov.com/routes"
	"github.com/gorilla/context"
	_ "github.com/mattn/go-sqlite3"
)

var err error

func main() {

	db.Init()

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
