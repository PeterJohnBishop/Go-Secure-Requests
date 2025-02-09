package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	// "net/http"

	"github.com/pquerna/otp/totp"
	// "github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

func hashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateSecretKey(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourAppName",
		AccountName: email,
	})
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}

	fmt.Println("Secret Key:", key.Secret())
	fmt.Println("QR Code URL:", key.URL())
	return key.Secret(), key.URL(), nil

}

// func generateQRCode(url string) {
// 	qrCodeURL := url
// 	err := qrcode.WriteFile(qrCodeURL, qrcode.Medium, 256, "totp_qr.png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("QR Code generated as totp_qr.png")
// }

// func generateQRCodeResponse(qrCodeURL string, w http.ResponseWriter, r *http.Request) {
// 	png, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
// 	if err != nil {
// 		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the response headers
// 	w.Header().Set("Content-Type", "image/png")
// 	w.Write(png)
// }

func verifyTOTP(userSecret string, otp string) {

	// User enters an OTP from their app
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otp)

	// Validate the OTP
	valid := totp.Validate(otp, userSecret)
	if valid {
		fmt.Println("OTP is valid!")
	} else {
		fmt.Println("Invalid OTP.")
	}
}
