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
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
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
