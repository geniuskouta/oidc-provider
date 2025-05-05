package main

import (
	"log"
	"net/http"
	"oidc/internal/oidc/handler"
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
	http.HandleFunc("/token", handler.ClientCredentialsFlow)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
