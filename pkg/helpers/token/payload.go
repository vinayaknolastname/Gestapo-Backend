package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Different types of error returend by token
var (
	ErrorExpiredToken error = fmt.Errorf("token is expired")
	ErrorInvalidToken error = fmt.Errorf("token is invalid")
)

const (
	sessionToken string = "session-token"
	accessToken  string = "access-token"
	ServiceToken string = "service-token"
)

// SessionPayload contains the payload data of the session token
type SessionPayload struct {
	Value     string `json:"value"`
	TokenFor  string `json:"token_for"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// AccessPayload contains the payload data of the token
type AccessPayload struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserType  string `json:"user_type"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// ServicePayload contains the payload data of the token and is used to validate between services
type ServicePayload struct {
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	ServiceName string `json:"service_name"`
	TokenType   string `json:"token_type"`
	jwt.RegisteredClaims
}

// NewSessionPayload creates a new token payload with a specific value and duration
func NewSessionPayload(value, tokenFor string) *SessionPayload {
	payload := &SessionPayload{
		Value:     value,
		TokenFor:  tokenFor,
		TokenType: sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload
}

// NewAccessPayload creates a new token payload with a specific username and duration
func NewAccessPayload(userID, userName, userType string) *AccessPayload {
	payload := &AccessPayload{
		UserID:    userID,
		UserName:  userName,
		UserType:  userType,
		TokenType: accessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)), //TODO:change token time
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload
}

// NewServicePayload creates a new token payload with a specific username and duration
func NewServicePayload(userID, userType, serviceName string) *ServicePayload {
	payload := &ServicePayload{
		UserID:      userID,
		UserType:    userType,
		ServiceName: serviceName,
		TokenType:   ServiceToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload
}

// Valid checks if the token payload is valid or not
func (payload *SessionPayload) Valid() error {
	if time.Now().After(payload.RegisteredClaims.ExpiresAt.Time) {
		return ErrorExpiredToken
	}
	return nil
}

func (payload *AccessPayload) Valid() error {
	if time.Now().After(payload.RegisteredClaims.ExpiresAt.Time) {
		return ErrorExpiredToken
	}
	return nil
}

func (payload *ServicePayload) Valid() error {
	if time.Now().After(payload.RegisteredClaims.ExpiresAt.Time) {
		return ErrorExpiredToken
	}
	return nil
}
