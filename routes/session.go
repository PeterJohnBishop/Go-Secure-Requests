package routes

// import (
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// )

// var ErrAuth = errors.New("authentication error")

// func PreAuthorize(r *http.Request) error {
// 	// Retrieve pending token from cookie
// 	t, err := r.Cookie("pending_2fa_token")
// 	if err != nil || t.Value == "" {
// 		return ErrAuth
// 	}

// 	pendingToken := t.Value

// 	// Find user by pending token
// 	var user *Login
// 	for _, u := range users {
// 		if u.Pending_2fa_Token == pendingToken {
// 			user = &u
// 			break
// 		}
// 	}

// 	if user == nil {
// 		return ErrAuth
// 	}

// 	return nil
// }

// func Authorize(r *http.Request) error {
// 	// Retrieve session token from cookie
// 	st, err := r.Cookie("session_token")
// 	if err != nil {
// 		fmt.Println("Error retrieving session_token:", err)
// 		return ErrAuth
// 	}

// 	if err != nil || st.Value == "" {
// 		fmt.Println("error verifying Session token.")
// 		return ErrAuth
// 	}

// 	sessionToken := st.Value

// 	// Find user by session token
// 	var user *Login
// 	for _, u := range users {
// 		if u.SessionToken == sessionToken {
// 			user = &u
// 			break
// 		}
// 	}

// 	if user == nil {
// 		return ErrAuth
// 	}

// 	// Retrieve and decode CSRF token
// 	csrf := r.Header.Get("X-CSRF-Token")
// 	fmt.Println("error verifying CSRF token.")
// 	decodedCSRF, err := url.QueryUnescape(csrf)
// 	if err != nil {
// 		fmt.Println("error verifying CSRF token.")
// 		return ErrAuth
// 	}

// 	// Compare decoded CSRF token with stored token
// 	if decodedCSRF == "" || decodedCSRF != user.CSRFToken {
// 		return ErrAuth
// 	}

// 	return nil
// }

// func Authorize(r *http.Request) error {
// 	// Log all cookies for debugging
// 	fmt.Println("Received Cookies:")
// 	for _, cookie := range r.Cookies() {
// 		fmt.Printf(" - %s: %s\n", cookie.Name, cookie.Value)
// 	}

// 	// Try to retrieve session token
// st, err := r.Cookie("session_token")
// if err != nil {
// 	fmt.Println("Error retrieving session_token:", err)
// 	return ErrAuth
// }

// 	if st.Value == "" {
// 		fmt.Println("Session token is empty")
// 		return ErrAuth
// 	}

// 	fmt.Println("Session token retrieved:", st.Value)

// 	// (Rest of the function remains unchanged)
// 	return nil
// }
