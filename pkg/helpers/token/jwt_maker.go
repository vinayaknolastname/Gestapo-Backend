package token

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JWTMaker is JSON Wed Token Maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateSessionToken implements Maker.
func (maker *JWTMaker) CreateSessionToken(value, tokenFor string) (string, error) {
	mySigningKey := []byte(maker.secretKey)
	payload := NewSessionPayload(value, tokenFor)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(mySigningKey)
}

// VerifySessionToken implements Maker.
func (maker *JWTMaker) VerifySessionToken(token string) (*SessionPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &SessionPayload{}, keyFunc)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}
	if !jwtToken.Valid {
		return nil, ErrorInvalidToken
	}
	payload, ok := jwtToken.Claims.(*SessionPayload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}

// CreateAccessToken create a token for specific userName and duration
func (maker *JWTMaker) CreateAccessToken(userID, userName, userType string) (string, error) {
	mySigningKey := []byte(maker.secretKey)
	payload := NewAccessPayload(userID, userName, userType)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(mySigningKey)

}

// VerifyAccessToken checks if token is valid or not
func (maker *JWTMaker) VerifyAccessToken(token string) (*AccessPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &AccessPayload{}, keyFunc)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}
	if !jwtToken.Valid {
		return nil, ErrorInvalidToken
	}
	payload, ok := jwtToken.Claims.(*AccessPayload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}

// CreateServiceToken create a token for specific service and duration
func (maker *JWTMaker) CreateServiceToken(userID, userType, serviceName string) (string, error) {
	mySigningKey := []byte(maker.secretKey)
	payload := NewServicePayload(userID, userType, serviceName)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(mySigningKey)
}

// VerifyServiceToken checks if token is valid or not
func (maker *JWTMaker) VerifyServiceToken(token string) (*ServicePayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &ServicePayload{}, keyFunc)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}
	if !jwtToken.Valid {
		return nil, ErrorInvalidToken
	}
	payload, ok := jwtToken.Claims.(*ServicePayload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}
