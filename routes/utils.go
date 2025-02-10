package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func GenerateSecretKey(email string) (string, string, error) {
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

func VerifyTOTP(userSecret string, otp string) bool {

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

func GenerateQRCodeBase64(textToEncode string) string {

	png, err := qrcode.Encode(textToEncode, qrcode.Medium, 256)
	if err != nil {
		return err.Error()
	}

	base64Image := base64.StdEncoding.EncodeToString(png)

	return base64Image
}

func GenerateQRCodePNG(textToEncode string) ([]byte, error) {

	png, err := qrcode.Encode(textToEncode, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return png, nil

}
