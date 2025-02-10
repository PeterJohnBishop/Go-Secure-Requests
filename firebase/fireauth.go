package firebase

import (
	"context"

	"firebase.google.com/go/auth"
)

func CreateUser(context context.Context, email string, password string) (string, bool) {
	user, err := authClient.CreateUser(context, (&auth.UserToCreate{}).
		Email("test@example.com").
		Password("securepassword"))
	if err != nil {
		return err.Error(), false
	}
	return user.UID, true
}
