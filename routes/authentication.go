package routes

import (
	"automatic-fiesta-go/main.go/firebase"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "time"
	"context"
)

func Register(w http.ResponseWriter, r *http.Request) {

	// step 1: create a user account and save the hashed password. Generate and save TOTP secret key.
	ctx := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(email) == 0 || len(password) == 0 {
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	uid, success := firebase.CreateUser(ctx, email, password)

	fmt.Println(uid)

	if success {
		secret, qrURL, err := GenerateSecretKey(email)
		if err != nil {
			log.Fatal("Error generating secret key:", err)
		}

		qrPng, error := GenerateQRCodePNG(qrURL)
		if error != nil {
			log.Fatal("Error generating QR Code:", error)
		}

		profileData := map[string]interface{}{
			"session_token": "",
			"csrf_token":    "",
			"temp_token":    "",
			"totp_secret":   secret,
		}

		firebase.CreateProfile(ctx, uid, profileData)

		w.Header().Set("Content-Type", "image/png")
		w.Write(qrPng)
	}
}

// func Login(w http.ResponseWriter, r *http.Request) {

// 	// step 2: Verify email / password. Generate a temp token on success.

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	contentType := r.Header.Get("Content-Type")
// 	if !strings.HasPrefix(contentType, "multipart/form-data") {
// 		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	email := r.FormValue("email")
// 	password := r.FormValue("password")

// 	user, ok := users[email]
// 	if !ok || !checkPasswordHash(password, user.HashedPassword) {
// 		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
// 		return
// 	}

// 	Pending_2fa_Token := generateToken(32)

// 	users[email] = Login{HashedPassword: user.HashedPassword, Pending_2fa_Token: Pending_2fa_Token, TOTPSecret: user.TOTPSecret}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "pending_2fa_token",
// 		Value:    Pending_2fa_Token,
// 		Expires:  time.Now().Add(5 * time.Minute), // 5min time limit
// 		HttpOnly: true,                            // true so the cookie is not accessible by the client
// 	})

// 	response := map[string]interface{}{
// 		"message": "Email and password validated. You have 5 minutes to complete TOTP Authentication.",
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func TOTP(w http.ResponseWriter, r *http.Request) {

// 	// step 3: Verify temp token, then verify TOTP code. On success generate and set session and crsf tokens. Clear temp token.

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	contentType := r.Header.Get("Content-Type")
// 	if !strings.HasPrefix(contentType, "multipart/form-data") {
// 		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	t, err := r.Cookie("pending_2fa_token")
// 	if err != nil || t.Value == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	pendingToken := t.Value

// 	errA := PreAuthorize(r)
// 	if errA != nil {
// 		fmt.Println(errA)
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var user *Login
// 	for _, u := range users {
// 		if u.Pending_2fa_Token == pendingToken {
// 			user = &u
// 			break
// 		}
// 	}

// 	if user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	user_otp := r.FormValue("otp")

// 	secondAuthPassed := verifyTOTP(user.TOTPSecret, user_otp)

// 	if !secondAuthPassed {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	sessionToken := generateToken(32)
// 	csrfToken := generateToken(32)

// 	// Set session token as cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session_token",
// 		Value:    sessionToken,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true, // true so the cookie is not accessible by the client
// 	})

// 	// Set CSRF token as cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "csrf_token",
// 		Value:    csrfToken,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: false, // false so the client can save and send it back for verification
// 	})

// 	// Store tokens in user object
// 	user.SessionToken = sessionToken
// 	user.CSRFToken = csrfToken
// 	user.Pending_2fa_Token = ""
// 	users["pjb.den@gmail.com"] = Login{SessionToken: sessionToken, CSRFToken: csrfToken, Pending_2fa_Token: ""} // hardcoded for testing!

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "pending_2fa_token",
// 		Value:   "",
// 		Expires: time.Now(),
// 	})

// 	response := map[string]interface{}{
// 		"message": "TOTP Authentication Successful.",
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func Protected(w http.ResponseWriter, r *http.Request) {

// 	// step 3: when a request is sent to the server the Authorize function verfies both tokens.

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	err := Authorize(r)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"message": "Protected route successfully accessed.",
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func Logout(w http.ResponseWriter, r *http.Request) {

// 	// step 4: on logout the session token and csrf token are revoked

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	contentType := r.Header.Get("Content-Type")
// 	if !strings.HasPrefix(contentType, "multipart/form-data") {
// 		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	err := Authorize(r)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	email := r.FormValue("email")
// 	user, ok := users[email]
// 	if !ok {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	user.SessionToken = ""
// 	user.CSRFToken = ""
// 	users[email] = user

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Value:   "",
// 		Expires: time.Now(),
// 	})

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "csrf_token",
// 		Value:   "",
// 		Expires: time.Now(),
// 	})

// 	response := map[string]interface{}{
// 		"message": "Logged Out",
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
