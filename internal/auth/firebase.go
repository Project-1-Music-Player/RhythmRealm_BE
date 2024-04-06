package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./credentials/rr-auth-467b4-firebase-adminsdk-1ortf-cc80f01084.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	AuthClient, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}
}
