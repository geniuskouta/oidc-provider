package service

import (
	"fmt"
	"oidc/internal/oidc/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey []byte
	issuer    string
	audience  string
}

func NewTokenService() (*TokenService, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return nil, fmt.Errorf("missing SECRET_KEY in environment")
	}

	issuer := os.Getenv("OIDC_ISSUER")
	if issuer == "" {
		issuer = "http://localhost:8080" // fallback or raise error
	}

	audience := os.Getenv("OIDC_AUDIENCE")
	if audience == "" {
		audience = "my-api"
	}

	return &TokenService{
		secretKey: []byte(secret),
		issuer:    issuer,
		audience:  audience,
	}, nil
}

func (s *TokenService) GenerateToken(subject string, scope string) (*domain.Token, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   subject,
		Audience:  []string{s.audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &domain.Token{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       scope,
	}, nil
}
