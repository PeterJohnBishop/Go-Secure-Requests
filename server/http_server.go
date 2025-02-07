package server

import (
	"fmt"
	"net/http"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

func Http_Server() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)
	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		http.Error(w, "Username or password is empty", http.StatusBadRequest)
		return
	}

	if _, ok := users[username]; ok {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := hashedPassword(password)
	users[username] = Login{HashedPassword: hashedPassword}

	fmt.Fprintf(w, "User %s has been registered", username)
}

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func protected(w http.ResponseWriter, r *http.Request) {}
