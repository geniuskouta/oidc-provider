package usecase

import (
	"crypto/subtle"
	"fmt"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/repo"
	"oidc/internal/oidc/service"
)

type ClientCredentialsFlow struct {
	repo    *repo.ClientRepository
	service *service.TokenService
}

func NewClientCredentialsFlow(repo *repo.ClientRepository, service *service.TokenService) *ClientCredentialsFlow {
	return &ClientCredentialsFlow{
		repo:    repo,
		service: service,
	}
}

func (u *ClientCredentialsFlow) Handle(clientID, clientSecret, grantType, scope string) (*domain.Token, error) {
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
	if !isValidClient(*u.repo, client) {
		return nil, fmt.Errorf("invalid_client")
	}

	// Generate the access token
	token, err := u.service.GenerateToken(clientID, scope)
	if err != nil {
		return nil, fmt.Errorf("server_error: %v", err)
	}

	return token, nil
}

func isValidClient(repo repo.ClientRepository, client *domain.Client) bool {
	target, err := repo.FindByClientID(client.ClientID)

	if err != nil || target == nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(target.ClientSecret)) == 1
}
