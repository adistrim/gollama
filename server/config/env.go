package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

var ENV *Config

func init() {
	var err error
	
	ENV, err = Load()
	if err != nil {
		log.Fatalf("Error: Failed to load environment variables: %v", err)
	}
}

func Load() (*Config, error) {
	godotenv.Load()
	
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("No PORT environment variable found, using default port 8080")
		port = "8080"
	}
	
	return &Config{
		Port: port,
	}, nil
}

