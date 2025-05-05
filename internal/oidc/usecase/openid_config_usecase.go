package usecase

import (
	"oidc/internal/oidc/domain"
	"os"
)

type OpenIDConfig struct {
	Address string
}

func NewOpenIDConfig() *OpenIDConfig {
	address := os.Getenv("OIDC_ADDRESS")
	if address == "" {
		address = "localhost:8080" // Default to localhost if not set
	}

	return &OpenIDConfig{
		Address: address,
	}
}

// GetConfig retrieves the OpenID configuration
func (u *OpenIDConfig) Get() *domain.OpenIDConfig {
	baseURL := "http://" + u.Address

	return &domain.OpenIDConfig{
		Issuer:                            baseURL,
		AuthorizationEndpoint:             baseURL + "/authorize",
		TokenEndpoint:                     baseURL + "/token",
		UserinfoEndpoint:                  baseURL + "/userinfo",
		JwksURI:                           baseURL + "/.well-known/jwks.json",
		RegistrationEndpoint:              baseURL + "/register",
		ScopesSupported:                   []string{"openid", "profile", "email"},
		ResponseTypesSupported:            []string{"code", "token", "id_token"},
		GrantTypesSupported:               []string{"authorization_code", "implicit", "client_credentials"},
		CodeChallengeMethodsSupported:     []string{"S256", "plain"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"},
	}
}
