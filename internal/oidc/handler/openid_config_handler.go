package handler

import (
	"encoding/json"
	"net/http"
	"oidc/internal/oidc/usecase"
)

type OpenIDConfig struct {
	usecase *usecase.OpenIDConfig
}

func NewOpenIDConfigHandler(usecase *usecase.OpenIDConfig) *OpenIDConfig {
	return &OpenIDConfig{
		usecase: usecase,
	}
}

func (h *OpenIDConfig) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	usecase := usecase.NewOpenIDConfig()
	config := usecase.Get()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "Unable to encode OpenID configuration", http.StatusInternalServerError)
	}
}
