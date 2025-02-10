package server

import (
	"automatic-fiesta-go/main.go/routes"
	"net/http"
)

type Login struct {
	HashedPassword    string
	SessionToken      string
	CSRFToken         string
	Pending_2fa_Token string
	TOTPSecret        string
}

var users = map[string]Login{}

func Http_Server() {
	mux := http.NewServeMux()

	mux.Handle("/register", http.HandlerFunc(routes.Register))
	// mux.Handle("/login", http.HandlerFunc(routes.Login))
	// mux.Handle("/2fa", http.HandlerFunc(routes.TOTP))
	// mux.Handle("/logout", http.HandlerFunc(routes.Logout))
	// mux.Handle("/protected", routes.SecureHeaders(routes.StrictSOPMiddleware(http.HandlerFunc(routes.Protected))))

	http.ListenAndServe(":8080", mux)
}
