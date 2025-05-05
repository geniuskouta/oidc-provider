// internal/oidc/domain/token.go
package domain

// Token represents the access token generated for the client credentials flow.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}
