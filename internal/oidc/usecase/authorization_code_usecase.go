package usecase

import (
	"crypto/subtle"
	"fmt"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/repo"
	"oidc/internal/oidc/service"
	"time"

	"github.com/google/uuid"
)

type AuthorizationCodeFlow struct {
	clientRepo   *repo.ClientRepository
	authCodeRepo *repo.AuthCodeRepository
	service      *service.TokenService
}

func NewAuthorizationCodeFlow(clientRepo *repo.ClientRepository, authCodeRepo *repo.AuthCodeRepository, service *service.TokenService) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{
		clientRepo:   clientRepo,
		authCodeRepo: authCodeRepo,
		service:      service,
	}
}

func (u *AuthorizationCodeFlow) Authorize(clientID, redirectURI, scope string) (string, error) {
	client, err := u.clientRepo.FindByClientID(clientID)
	if err != nil || client == nil {
		return "", fmt.Errorf("unauthorized_client")
	}

	// TODO: Validate redirectURI against registered ones
	if !isValidRedirectUri(client, redirectURI) {
		return "", fmt.Errorf("invalid_redirect_uri")
	}

	code := uuid.NewString()
	authCode := &domain.AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	err = u.authCodeRepo.Save(authCode)
	if err != nil {
		return "", fmt.Errorf("failed to save authorization code: %w", err)
	}

	return code, nil
}

func (u *AuthorizationCodeFlow) Exchange(code, clientID string) (*domain.Token, error) {
	if !isValidAuthCode(*u.authCodeRepo, code) {
		return nil, fmt.Errorf("invalid_auth_code")
	}

	// Generate token using token service
	token, err := u.service.GenerateToken(clientID, "openid")
	if err != nil {
		return nil, err
	}

	return token, nil
}

func isValidAuthCode(repo repo.AuthCodeRepository, code string) bool {
	target, err := repo.FindByCode(code)

	if err != nil || target == nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(code), []byte(target.Code)) == 1
}

func isValidRedirectUri(client *domain.Client, redirectURI string) bool {
	for _, uri := range client.RedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}
