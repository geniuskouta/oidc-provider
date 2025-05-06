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
	tokenService, err := service.NewTokenService()

	if err != nil {
		log.Fatalf("failed to initialize token service: %v", err)
	}

	clientCredentialsFlowUsecase := usecase.NewClientCredentialsFlow(clientRepo, tokenService)
	registerClientUsecase := usecase.NewRegisterClient(clientRepo)
	openIDConfigUsecase := usecase.NewOpenIDConfig()

	clientCredentialsFlowHandler := handler.NewClientCredentialsFlow(clientCredentialsFlowUsecase)
	registerClientHandler := handler.NewRegisterClient(registerClientUsecase)
	openIDConfigHandler := handler.NewOpenIDConfigHandler(openIDConfigUsecase)

	r := mux.NewRouter()

	r.HandleFunc("/token", clientCredentialsFlowHandler.Handle).Methods("POST")
	r.HandleFunc("/register", registerClientHandler.Handle).Methods("POST")
	r.HandleFunc("/.well-known/openid-configuration", openIDConfigHandler.Handle).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("OIDC_ADDRESS"), r))
}
