package auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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

	res.Write([]byte("Hello" + databaseUsername))

}
