package handler

import (
	"encoding/json"
	"net/http"
	"oidc/internal/oidc/usecase"
)

type AuthorizationCodeFlow struct {
	usecase *usecase.AuthorizationCodeFlow
}

func NewAuthorizationCodeFlow(usecase *usecase.AuthorizationCodeFlow) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{usecase: usecase}
}

func (h *AuthorizationCodeFlow) HandleAuthorizationCode(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")

	code, err := h.usecase.Authorize(clientID, redirectURI, scope)
	if err != nil {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirectURI+"?code="+code, http.StatusFound)
}

func (h *AuthorizationCodeFlow) HandleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	code := r.PostFormValue("code")
	clientID := r.PostFormValue("client_id")

	token, err := h.usecase.Exchange(code, clientID)
	if err != nil {
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
