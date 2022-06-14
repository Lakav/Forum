package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
