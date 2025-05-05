package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/repo"
	"time"
)

type RegisterClient struct {
	repo *repo.ClientRepository
}

func NewRegisterClient(repo *repo.ClientRepository) *RegisterClient {
	return &RegisterClient{
		repo: repo,
	}
}

func (u *RegisterClient) Handle(appName string, redirectURIs []string) (*domain.Client, error) {
	if appName == "" || len(redirectURIs) == 0 {
		return nil, errors.New("invalid_request: app name and redirect URIs required")
	}

	clientID := generateRandomID()
	clientSecret := generateRandomID()

	client := &domain.Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         appName,
		RedirectURIs: redirectURIs,
		CreatedAt:    time.Now().Unix(),
	}

	if err := u.repo.Save(client); err != nil {
		return nil, err
	}

	return client, nil
}

func generateRandomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
