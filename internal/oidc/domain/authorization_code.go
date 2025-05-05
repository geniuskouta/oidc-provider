package domain

import "time"

type AuthorizationCode struct {
	Code        string
	ClientID    string
	RedirectURI string
	Scope       string
	ExpiresAt   time.Time
}
