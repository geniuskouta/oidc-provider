package repo

import (
	"errors"
	"fmt"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/infra"
)

type ClientRepository struct {
	db *infra.DB
}

func NewClientRepo(db *infra.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (repo *ClientRepository) Save(client *domain.Client) error {
	// Check if client already exists (based on client ID)
	if _, err := repo.db.Get(client.ClientID); err == nil {
		return fmt.Errorf("client with ID %s already exists", client.ClientID)
	}

	// Save the client to the database
	repo.db.Set(client.ClientID, client)
	return nil
}

func (repo *ClientRepository) FindByClientID(clientID string) (*domain.Client, error) {
	// Fetch the client from the database
	clientData, err := repo.db.Get(clientID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	// Type assert and return the client
	client, ok := clientData.(*domain.Client)
	if !ok {
		return nil, errors.New("failed to cast client data")
	}

	return client, nil
}
