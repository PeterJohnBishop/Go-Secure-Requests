package routes

import (
	"automatic-fiesta-go/main.go/firebase"
	"encoding/json"
	"time"

	// "encoding/json"

	"net/http"
	"strings"

	"context"
)

func Register(w http.ResponseWriter, r *http.Request) {

	// Send email and password to create a user in Firebase Authentication.
	// On success generate a TOTP Secret Key, generate TOTP QR code, save secret key to new Profile doc in Firestore
	// return QR code image

	ctx := context.Background()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	if len(email) == 0 || len(password) == 0 {
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	uid, success := firebase.CreateUser(ctx, email, password)
	if !success {
		http.Error(w, "Error creating Firebase User", http.StatusBadRequest)
		return
	}

	secret, qrURL, err := GenerateSecretKey(email)
	if err != nil {
		http.Error(w, "Error creating secret key", http.StatusBadRequest)
		return
	}

	qrPng, err := GenerateQRCodePNG(qrURL)
	if err != nil {
		http.Error(w, "Error creating QR Code image", http.StatusBadRequest)
		return
	}

	profileData := map[string]interface{}{
		"session_token": "",
		"csrf_token":    "",
		"temp_token":    "",
		"totp_secret":   secret,
	}

	_, created := firebase.CreateProfile(ctx, uid, profileData)
	if !created {
		http.Error(w, "Error creating Firestore Profile Doc", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qrPng)

}

func Verify(w http.ResponseWriter, r *http.Request) {

	// After Client side login verify Firebase UserIDToken generate and save a temp token
	// Set temp token as cookie and send user for TOTP authentication

	ctx := context.Background()

	uid := r.FormValue("uid")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	success := firebase.VerifyIDToken(ctx, token)
	if !success {
		http.Error(w, "Token verification failed.", http.StatusUnauthorized)
		return
	}

	temp := GenerateToken(32)
	_, updated := firebase.UpdateProfileField(ctx, uid, "temp_token", temp)
	if !updated {
		http.Error(w, "Error saving temp token.", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "temp",
		Value:    temp,
		Expires:  time.Now().Add(5 * time.Minute), // 5min time limit
		HttpOnly: true,                            // true so the cookie is not accessible by the client
	})

	response := map[string]interface{}{
		"message": "Email and password validated. You have 5 minutes to complete TOTP Authentication.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func TOTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	uid := r.FormValue("uid")
	user_otp := r.FormValue("otp")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Invalid Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	userProfile, _, found := firebase.GetProfile(ctx, uid)
	if !found {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	preAuth := PreAuthorize(ctx, userProfile, r)
	if preAuth != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	secondAuthPassed := VerifyTOTP(userProfile.TOTPSecret, user_otp)

	if !secondAuthPassed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := GenerateToken(32)
	csrfToken := GenerateToken(32)

	// Set session token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // true so the cookie is not accessible by the client
	})

	// Set CSRF token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false, // false so the client can save and send it back for verification
	})

	updates := map[string]interface{}{
		"session_token": sessionToken,
		"csrf_token":    csrfToken,
		"temp_token":    "",
	}

	_, updated := firebase.UpdateMultipleProfileFields(ctx, uid, updates)
	if !updated {
		http.Error(w, "Updating profile tokens failed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "pending_2fa_token",
		Value:   "",
		Expires: time.Now(),
	})

	response := map[string]interface{}{
		"message": "TOTP Authentication Successful.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
