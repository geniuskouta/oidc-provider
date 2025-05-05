// handler/register_client.go
package handler

import (
	"encoding/json"
	"net/http"
	"oidc/internal/oidc/usecase"
)

type RegisterClient struct {
	usecase *usecase.RegisterClient
}

func NewRegisterClient(usecase *usecase.RegisterClient) *RegisterClient {
	return &RegisterClient{
		usecase: usecase,
	}
}

func (h *RegisterClient) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientName   string   `json:"client_name"`
		RedirectURIs []string `json:"redirect_uris"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	client, err := h.usecase.Handle(req.ClientName, req.RedirectURIs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
