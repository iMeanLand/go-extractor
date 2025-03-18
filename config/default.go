package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	ApiKey = ""
	ApiUrl = ""
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ApiKey = os.Getenv("API_KEY")
	ApiUrl = os.Getenv("API_URL")

	if ApiKey == "" || ApiUrl == "" {
		log.Fatal("API_KEY or API_URL env variables not set")
	}
}
