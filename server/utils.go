package server

import (
	"golang.org/x/crypto/bcrypt"
)

func hashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}
