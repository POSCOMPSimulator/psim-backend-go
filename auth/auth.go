package auth

import (
	"context"
	"errors"
	"os"

	"google.golang.org/api/idtoken"
	"poscomp-simulator.com/backend/models"
)

func VerifyIdToken(idToken string) (models.Usuario, error) {

	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("CLIENT_ID"))
	if err != nil {
		return models.Usuario{}, err
	}

	if payload.Audience != os.Getenv("CLIENT_ID") {
		return models.Usuario{}, errors.New("Token inválido.")
	}

	u := models.Usuario{
		Email:      payload.Claims["email"].(string),
		Nome:       payload.Claims["name"].(string),
		FotoPerfil: payload.Claims["picture"].(string),
	}
	return u, nil

}
