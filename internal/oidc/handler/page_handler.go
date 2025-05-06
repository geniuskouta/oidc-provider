package handler

import (
	"net/http"
	"oidc/internal/oidc/templates"
)

// PageHandlers struct to group page-related handlers
type PageHandler struct {
	templateEngine *templates.Engine
}

func NewPageHandler(templateEngine *templates.Engine) *PageHandler {
	return &PageHandler{
		templateEngine: templateEngine,
	}
}

// RenderLoginPage renders the login page using the template engine
func (h *PageHandler) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for the template (can also get from query params or context)
	data := map[string]interface{}{
		"Email":       "", // Empty initial value for email
		"ClientID":    r.URL.Query().Get("client_id"),
		"RedirectURI": r.URL.Query().Get("redirect_uri"),
		"Scope":       r.URL.Query().Get("scope"),
	}

	// Render the login page
	if err := h.templateEngine.Render(w, "login.html", data); err != nil {
		http.Error(w, "Failed to render login page", http.StatusInternalServerError)
	}
}
