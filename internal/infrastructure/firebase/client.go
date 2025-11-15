package firebase

import (
	"context"
	"log"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client holds the initialized Firebase Auth client.
type Client struct {
	AuthClient *auth.Client
}

// NewClient initializes the Firebase Admin SDK using the service account file.
func NewClient(serviceAccountKeyPath string) *Client {
	ctx := context.Background()

	// 1. Create options using the path to the service account JSON
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	
	// 2. Initialize the Firebase App
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	// 3. Get the Authentication client instance
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}

	return &Client{
		AuthClient: authClient,
	}
}