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

func CreateProfile(ctx context.Context, uid string, profile map[string]interface{}) (string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	_, err := firestoreClient.Collection("profiles").Doc(uid).Set(ctx, profile)
	if err != nil {
		return err.Error(), false
	}
	fmt.Println("Profile added successfully!")
	return "Success", true
}

func GetProfile(ctx context.Context, uid string) (Profile, string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	docRef := firestoreClient.Collection("profiles").Doc(uid)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return Profile{}, err.Error(), false
	}
	var profile Profile
	if err := docSnap.DataTo(&profile); err != nil {
		return Profile{}, err.Error(), false
	}
	return profile, "Success", true

}

func UpdateProfileField(ctx context.Context, uid string, field string, newValue interface{}) (string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	docRef := firestoreClient.Collection("profiles").Doc(uid)

	// Update the specific field
	_, err := docRef.Update(ctx, []firestore.Update{
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

func UpdateMultipleProfileFields(ctx context.Context, uid string, updates map[string]interface{}) (string, bool) {

	// Create a new context with a 60-second timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	docRef := firestoreClient.Collection("profiles").Doc(uid)

	var firestoreUpdates []firestore.Update
	for key, value := range updates {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	_, err := docRef.Update(ctx, firestoreUpdates)
	if err != nil {
		return err.Error(), false
	}

	return "Success", true
}
