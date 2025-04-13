package mattercontroller

import (
	"context"
	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func initFirebase() *firestore.Client {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	serviceAccountKeyPath := os.Getenv("FIREBASE_SERVER_CREDENTIALS")

	// Initialize Firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
