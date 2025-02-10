package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

var authClient *auth.Client
var firestoreClient *firestore.Client

func Init() error {
	opt := option.WithCredentialsFile("firebase/automatic-fiesta-4fe57-firebase-adminsdk-fbsvc-923b968f5c.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Auth client: %v\n", err)
	}

	firestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}

	return nil
}
