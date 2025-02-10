package firebase

import (
	"context"
	"fmt"
	"time"
)

type Profile struct {
	SessionToken string
	CSRFToken    string
	TempToken    string
	TOTPSecret   string
}

func CreateProfile(ctx context.Context, documentId string, profile map[string]interface{}) error {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	_, err := firestoreClient.Collection("profiles").Doc(documentId).Set(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to add item: %v", err)
	}
	fmt.Println("Item added successfully!")
	return nil
}
