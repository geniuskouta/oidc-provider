package repo

import (
	"fmt"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/infra"
	"time"
)

type AuthCodeRepository struct {
	db *infra.DB
}

func NewAuthCodeRepo(db *infra.DB) *AuthCodeRepository {
	return &AuthCodeRepository{
		db: db,
	}
}

func (r *AuthCodeRepository) Save(authCode *domain.AuthorizationCode) error {
	r.db.Set(authCode.Code, authCode)
	return nil
}

func (r *AuthCodeRepository) FindByCode(code string) (*domain.AuthorizationCode, error) {
	value, err := r.db.Get(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find authorization code: %w", err)
	}

	authCode, ok := value.(*domain.AuthorizationCode)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for authorization code")
	}

	if time.Now().After(authCode.ExpiresAt) {
		return nil, fmt.Errorf("authorization code has expired")
	}

	return authCode, nil
}
