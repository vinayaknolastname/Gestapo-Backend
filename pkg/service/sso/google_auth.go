package sso

import (
	"context"
	"fmt"

	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"google.golang.org/api/idtoken"
)

const missingClaims string = "missing claims"

func GoogleOauth(token, clientID string, log logger.Logger) (string, string, error) {
	idtoken, err := idtoken.Validate(context.Background(), token, clientID)
	if err != nil {
		log.LogError("error while validating idtoken in GoogleAndroidOauth:", err.Error())
		return "", "", err
	}
	if idtoken.Claims == nil {
		log.LogError(missingClaims)
		return "", "", fmt.Errorf(missingClaims)
	}
	email, emailExist := idtoken.Claims["email"].(string)
	if !emailExist {
		log.LogError("missing email in calims")
		return "", "", fmt.Errorf(missingClaims)
	}
	name, nameExist := idtoken.Claims["name"].(string)
	if !nameExist {
		log.LogError("missing name in calims")
		return "", "", fmt.Errorf(missingClaims)
	}
	return email, name, nil

}
