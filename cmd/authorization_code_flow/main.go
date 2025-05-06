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

	authorizationCodeFlowHandler := handler.NewAuthorizationCodeFlow(authorizationCodeFlowUsecase)
	registerClientHandler := handler.NewRegisterClient(registerClientUsecase)
	openIDConfigHandler := handler.NewOpenIDConfigHandler(openIDConfigUsecase)

	http.HandleFunc("/signup", authorizationCodeFlowHandler.HandleSignUpUser)
	http.HandleFunc("/authorize", authorizationCodeFlowHandler.StartAuthorization)
	http.HandleFunc("/login", authorizationCodeFlowHandler.HandleAuthorizationCode)
	http.HandleFunc("/token", authorizationCodeFlowHandler.HandleToken)
	http.HandleFunc("/register", registerClientHandler.Handle)
	http.HandleFunc("/.well-known/openid-configuration", openIDConfigHandler.Handle)
	log.Fatal(http.ListenAndServe(os.Getenv("OIDC_ADDRESS"), nil))
}
