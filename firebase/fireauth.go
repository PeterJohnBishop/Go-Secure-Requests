package firebase

import (
	"context"

	"firebase.google.com/go/auth"
)

func CreateUser(context context.Context, email string, password string) (string, bool) {
	user, err := authClient.CreateUser(context, (&auth.UserToCreate{}).
		Email(email).
		Password(password))
	if err != nil {
		return err.Error(), false
	}
	return user.UID, true
}

func VerifyIDToken(context context.Context, token string) bool {
	_, err := authClient.VerifyIDToken(context, token)
	return err == nil
}
