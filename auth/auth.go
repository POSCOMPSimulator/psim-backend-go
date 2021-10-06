package auth

import (
	"net/http"

	"google.golang.org/api/oauth2/v2"
)

func VerifyIdToken(idToken string) (string, error) {

	var httpClient = &http.Client{}

	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()

	if err != nil {
		return "", err
	}

	return tokenInfo.UserId, nil
}
