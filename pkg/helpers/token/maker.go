package token

// Maker is an interface for managing tokens
type Maker interface {
	// CreateSessionToken create a session token for specific value and duration
	CreateSessionToken(value, tokenFor string) (string, error)

	// VerifySessionToken  checks if session token is valid or not
	VerifySessionToken(token string) (*SessionPayload, error)

	// CreateAccessToken create a access token for specific userName and duration
	CreateAccessToken(userId, userName, userType string) (string, error)

	// VerifyAccessToken checks if access token is valid or not
	VerifyAccessToken(token string) (*AccessPayload, error)

	// CreateServiceToken create a service token
	CreateServiceToken(userID, userType, serviceName string) (string, error)

	// VerifyServiceToken  checks if service token is valid or not
	VerifyServiceToken(token string) (*ServicePayload, error)
}
