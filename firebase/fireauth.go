package firebase

import (
	"context"
	"time"

	"firebase.google.com/go/auth"
)

func CreateUser(ctx context.Context, email string, password string) (string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	user, err := authClient.CreateUser(ctx, (&auth.UserToCreate{}).
		Email(email).
		Password(password))
	if err != nil {
		return err.Error(), false
	}
	return user.UID, true
}

func VerifyIDToken(ctx context.Context, token string) bool {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	_, err := authClient.VerifyIDToken(ctx, token)
	return err == nil
}
