package auth

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

func VerifyIdToken(idToken string) (map[string]interface{}, bool, error) {

	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("CLIENT_ID"))
	if err != nil {
		panic(err)
	}

	if payload.Audience != os.Getenv("CLIENT_ID") {
		return nil, false, nil
	}

	return payload.Claims, true, nil
}
