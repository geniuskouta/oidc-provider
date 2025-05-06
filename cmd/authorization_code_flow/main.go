package main

import (
	"log"
	"net/http"
	"oidc/internal/oidc/handler"
	"oidc/internal/oidc/infra"
	"oidc/internal/oidc/repo"
	"oidc/internal/oidc/service"
	"oidc/internal/oidc/usecase"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	address := os.Getenv("OIDC_ADDRESS")
	secretKey := os.Getenv("SECRET_KEY")

	if address == "" || secretKey == "" {
		log.Fatal("Missing required environment variables")
	}
}

func main() {
	db := infra.NewDB()
	clientRepo := repo.NewClientRepo(db)
	userRepo := repo.NewUserRepo(db)
	authCodeRepo := repo.NewAuthCodeRepo(db)
	tokenService, err := service.NewTokenService()

	if err != nil {
		log.Fatalf("failed to initialize token service: %v", err)
	}

	authorizationCodeFlowUsecase := usecase.NewAuthorizationCodeFlow(
		clientRepo,
		authCodeRepo,
		userRepo,
		tokenService,
	)
	registerClientUsecase := usecase.NewRegisterClient(clientRepo)
	openIDConfigUsecase := usecase.NewOpenIDConfig()

	// Create the handlers
	authorizationCodeFlowHandler := handler.NewAuthorizationCodeFlow(authorizationCodeFlowUsecase)
	registerClientHandler := handler.NewRegisterClient(registerClientUsecase)
	openIDConfigHandler := handler.NewOpenIDConfigHandler(openIDConfigUsecase)

	// Create a new Mux router
	r := mux.NewRouter()

	// Define the routes and associate handlers
	r.HandleFunc("/signup", authorizationCodeFlowHandler.HandleSignUpUser).Methods("POST")
	r.HandleFunc("/authorize", authorizationCodeFlowHandler.StartAuthorization).Methods("GET")
	r.HandleFunc("/login", authorizationCodeFlowHandler.HandleAuthorizationCode).Methods("POST")
	r.HandleFunc("/token", authorizationCodeFlowHandler.HandleToken).Methods("POST")
	r.HandleFunc("/register", registerClientHandler.Handle).Methods("POST")
	r.HandleFunc("/.well-known/openid-configuration", openIDConfigHandler.Handle).Methods("GET")

	// Start the HTTP server with the Mux router
	log.Fatal(http.ListenAndServe(os.Getenv("OIDC_ADDRESS"), r))
}
