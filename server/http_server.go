package server

import (
	"fmt"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

// Cross-Site Request Forgery (CSRF): now that the above session
// token is being sent with every request, if a malicious site
// triggers a request from my machine, it will contain the session
// token. Allowing any request sent from that site to be authenticated
// as valid. To prevent this, a CSRF token is generated and sent to the
// client. This token is then sent back with every request. The server
// can then verify that the token is correct and that the request is
// not a CSRF attack.

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

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// Set session token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	// Set CSRF token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	// Store tokens in user object
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	fmt.Fprintln(w, "Login successful")
}

func logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := Authorize(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	user, ok := users[username]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.SessionToken = ""
	user.CSRFToken = ""
	users[username] = user

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "csrf_token",
		Value:   "",
		Expires: time.Now(),
	})

	fmt.Fprintln(w, "Logout successful")
}

func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := Authorize(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Protected resource")
}
