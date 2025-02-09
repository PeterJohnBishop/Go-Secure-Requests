package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/pquerna/otp/totp"
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

func verifyTOTP(userSecret string, otp string) bool {

	// User enters an OTP from their app
	fmt.Print("Enter OTP: ")
	fmt.Scanln(&otp)

	// Validate the OTP
	valid := totp.Validate(otp, userSecret)
	if valid {
		fmt.Println("OTP is valid!")
		return true
	} else {
		fmt.Println("Invalid OTP.")
		return false
	}
}
