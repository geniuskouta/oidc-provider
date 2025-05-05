package usecase

import (
	"fmt"
	"oidc/internal/oidc/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ClientCredentialsFlow(clientID, clientSecret, grantType, scope string) (*domain.Token, error) {
	// Validate grant_type
	if grantType != "client_credentials" {
		return nil, fmt.Errorf("unsupported_grant_type")
	}

	// Create client object
	client := &domain.Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Validate client credentials here in the usecase layer
	if !isValidClient(client) {
		return nil, fmt.Errorf("invalid_client")
	}

	// Load the secret key from environment variable
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if secretKey == nil || len(secretKey) == 0 {
		return nil, fmt.Errorf("secret key is missing")
	}

	// Generate the access token
	accessToken, err := generateAccessToken(clientID, scope, secretKey)
	if err != nil {
		return nil, fmt.Errorf("server_error: %v", err)
	}

	// Return the token response
	token := &domain.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       scope,
	}

	return token, nil
}

func isValidClient(client *domain.Client) bool {
	return client.ClientID == "my-client-id" && client.ClientSecret == "my-client-secret"
}

func generateAccessToken(clientID, scope string, secretKey []byte) (string, error) {
	// Create JWT claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "my-auth-server",                                  // Issuer identifier (Auth server)
		Subject:   clientID,                                          // Subject (the client ID)
		Audience:  []string{"my-api"},                                // Audience (the resource server)
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 hour expiry
		IssuedAt:  jwt.NewNumericDate(time.Now()),                    // Issued at
	}

	// Create a new token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
