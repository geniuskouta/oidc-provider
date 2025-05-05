package main

import (
	"log"
	"net/http"
	"oidc/internal/oidc/handler"
	"oidc/internal/oidc/infra"
	"oidc/internal/oidc/repo"
	"oidc/internal/oidc/usecase"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is required")
	}

	db := infra.NewDB()
	clientRepo := repo.NewClientRepo(db)

	clientCredentialsFlowUsecase := usecase.NewClientCredentialsFlow(clientRepo)
	registerClientUsecase := usecase.NewRegisterClient(clientRepo)

	clientCredentialsFlowHandler := handler.NewClientCredentialsFlow(clientCredentialsFlowUsecase)
	registerClientHandler := handler.NewRegisterClient(registerClientUsecase)

	http.HandleFunc("/token", clientCredentialsFlowHandler.Handle)
	http.HandleFunc("/register", registerClientHandler.Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
