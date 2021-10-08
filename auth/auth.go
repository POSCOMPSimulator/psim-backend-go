package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"google.golang.org/api/idtoken"
	"poscomp-simulator.com/backend/models"
)

func VerifyIdToken(auth string) (models.Usuario, error) {

	authorization := strings.Fields(auth)

	if len(authorization) != 2 || authorization[0] != "Bearer" {
		return models.Usuario{}, errors.New("Formato de autenticação incorreto.")
	}

	idToken := authorization[1]

	if idToken == os.Getenv("DUMMY_TOKEN") {
		u := models.Usuario{}
		u.GetDummy()
		return u, nil
	}

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
