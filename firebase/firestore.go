package firebase

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type Profile struct {
	SessionToken string
	CSRFToken    string
	TempToken    string
	TOTPSecret   string
}

func CreateProfile(ctx context.Context, documentId string, profile map[string]interface{}) (string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	_, err := firestoreClient.Collection("profiles").Doc(documentId).Set(ctx, profile)
	if err != nil {
		return err.Error(), false
	}
	fmt.Println("Profile added successfully!")
	return "Success", true
}

func UpdateProfileField(userID string, field string, newValue interface{}) (string, bool) {
	docRef := firestoreClient.Collection("profiles").Doc(userID)

	// Update the specific field
	_, err := docRef.Update(context.Background(), []firestore.Update{
		{
			Path:  field,
			Value: newValue,
		},
	})
	if err != nil {
		return err.Error(), false
	}
	return "Success", true
}
