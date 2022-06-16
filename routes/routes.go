package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"text/template"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var tpl *template.Template

var store = sessions.NewCookieStore([]byte("super-secret"))

func LoginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	row := db.QueryRow("SELECT username, password FROM users WHERE username=\"?\" LIMIT 1;", username)

	err := row.Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	if err == nil {
		session, _ := store.Get(r, "session")
		session.Values["userID"] = userID
		session.Save(r, w)
		tpl.ExecuteTemplate(w, "index.html", "Logged In")

	}
	fmt.Println("incorect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****indexHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok: ", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.statusfound is 302
		return
	}
	tpl.ExecuteTemplate(w, "index.html", "Logged In")
}

func SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}
