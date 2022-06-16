package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"text/template"

	"forumynov.com/db"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

func LoginPage(res http.ResponseWriter, req *http.Request) {
	tpl := template.Must(template.ParseFiles("public/login.html"))
	if req.Method != "POST" {
		tpl.Execute(res, struct {
			Error string
		}{
			Error: "",
		})
		return
	}
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	fmt.Println(username, password)
	rows, err := db.DB.Query("SELECT username, password FROM users WHERE username=? LIMIT 1;", username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(res, req, "/login", 301)
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&databaseUsername, &databasePassword)
	if err != nil {
		fmt.Println(err)
		tpl.Execute(res, struct {
			Error string
		}{
			Error: "check username",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))

	if err != nil {
		fmt.Println(err)
		fmt.Println("incorect password")
		tpl.Execute(res, struct {
			Error string
		}{
			Error: "check username and password",
		})
		return
	}
	session, _ := store.Get(req, "session")
	var user string

	userID := db.DB.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)
	session.Values["id"] = userID
	session.Save(req, res)

	http.Redirect(res, req, "/", 301)

}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("public/index.html"))
	fmt.Println("*****indexHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["id"]
	fmt.Println("ok: ", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.statusfound is 302
		return
	}
	tpl.Execute(w, "Logged In")
}

func SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "public/signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.DB.QueryRow("SELECT username FROM users WHERE username=\"?\" LIMIT 1", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.DB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, string(hashedPassword))
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
