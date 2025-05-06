package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"oidc/internal/oidc/usecase"
)

type AuthorizationCodeFlow struct {
	usecase *usecase.AuthorizationCodeFlow
}

func NewAuthorizationCodeFlow(usecase *usecase.AuthorizationCodeFlow) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{usecase: usecase}
}

func (h *AuthorizationCodeFlow) HandleSignUpUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Println("request")

	// Call the SignUp use case
	err := h.usecase.SignUp(req.Email, req.Password)

	if err != nil {
		// Check for specific errors and handle them accordingly
		if err.Error() == "user already exists" {
			http.Error(w, "user already exists", http.StatusConflict) // 409 Conflict for existing user
		} else {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity) // 422 for other validation errors
		}
		return
	}

	// Successful sign-up, send response
	w.WriteHeader(http.StatusCreated) // 201 Created
}

func (h *AuthorizationCodeFlow) StartAuthorization(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")

	loginUrl, err := h.usecase.GetLoginUrl(clientID, redirectURI, scope)
	if err != nil {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, loginUrl, http.StatusFound)
}

func (h *AuthorizationCodeFlow) HandleAuthorizationCode(w http.ResponseWriter, r *http.Request) {
	// Parse form input (POSTed by the login form)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	scope := r.FormValue("scope")

	// 1. Authenticate user
	valid, err := h.usecase.AuthenticateUser(email, password)
	if err != nil || !valid {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 2. Authorize and generate code
	code, err := h.usecase.Authorize(clientID, redirectURI, scope)
	if err != nil {
		http.Error(w, "Authorization failed", http.StatusBadRequest)
		return
	}

	// 3. Redirect back to client
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
