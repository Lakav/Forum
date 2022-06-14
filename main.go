package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"time"
)

var users = map[string]string{"user1": "password", "user2": "password"}

var store = sessions.NewCookieStore([]byte("my_secret_key"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	r.HandleFunc("/logintest", logintest).Methods("GET")

	httpServer := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	rows.Close() //good habit to close

	// delete
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}

func logintest(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./login.html"))
	_ = tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Check if user exists
	storedPassword, exists := users[username]
	if exists {
		// It returns a new session if the sessions doesn't exist
		session, _ := store.Get(r, "session.id")
		if storedPassword == password {
			session.Values["authenticated"] = true
			// Saves all sessions used during the current request
			session.Save(r, w)
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		}
		w.Write([]byte("Login successfully!"))
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get registers and returns a session for the given name and session store.
	session, _ := store.Get(r, "session.id")
	// Set the authenticated value on the session to false
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Logout Successful"))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	if authenticated != nil && authenticated != false {
		w.Write([]byte("Welcome!"))
		return
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
