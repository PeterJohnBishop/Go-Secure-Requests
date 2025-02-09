package server

import (
	"errors"
	"net/http"
	"net/url"
)

var ErrAuth = errors.New("authentication error")

func PreAuthorize(r *http.Request) error {
	// Retrieve pending token from cookie
	t, err := r.Cookie("pending_2fa_token")
	if err != nil || t.Value == "" {
		return ErrAuth
	}

	pendingToken := t.Value

	// Find user by pending token
	var user *Login
	for _, u := range users {
		if u.Pending_2fa_Token == pendingToken {
			user = &u
			break
		}
	}

	if user == nil {
		return ErrAuth
	}

	return nil
}

func Authorize(r *http.Request) error {
	// Retrieve session token from cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return ErrAuth
	}

	sessionToken := st.Value

	// Find user by session token
	var user *Login
	for _, u := range users {
		if u.SessionToken == sessionToken {
			user = &u
			break
		}
	}

	if user == nil {
		return ErrAuth
	}

	// Retrieve and decode CSRF token
	csrf := r.Header.Get("X-CSRF-Token")
	decodedCSRF, err := url.QueryUnescape(csrf)
	if err != nil {
		return ErrAuth
	}

	// Compare decoded CSRF token with stored token
	if decodedCSRF == "" || decodedCSRF != user.CSRFToken {
		return ErrAuth
	}

	return nil
}
