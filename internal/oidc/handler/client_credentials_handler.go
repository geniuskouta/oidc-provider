// internal/oidc/handler/token_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"oidc/internal/oidc/usecase"
)

type ClientCredentialsFlow struct {
	usecase *usecase.ClientCredentialsFlow
}

func NewClientCredentialsFlow(usecase *usecase.ClientCredentialsFlow) *ClientCredentialsFlow {
	return &ClientCredentialsFlow{
		usecase: usecase,
	}
}

func (h *ClientCredentialsFlow) Handle(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract client credentials from Basic Auth header
	clientID, clientSecret, err := extractClientCredentials(r)
	if err != nil {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	// Parse the JSON request body
	var req struct {
		GrantType string `json:"grant_type"`
		Scope     string `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Handle(clientID, clientSecret, req.GrantType, req.Scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// Helper function to extract client credentials from Basic Auth header
func extractClientCredentials(r *http.Request) (string, string, error) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return "", "", fmt.Errorf("missing client credentials in Authorization header")
	}
	return clientID, clientSecret, nil
}
