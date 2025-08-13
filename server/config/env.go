package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	BaseURL string
	GithubToken string
	DatabaseURL string
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
	
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Println("No BASE_URL environment variable found, using default base URL http://localhost:11434/v1")
		baseURL = "http://localhost:11434/v1"
	}
	
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Println("No GITHUB_TOKEN environment variable found")
	}
	
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("No DATABASE_URL environment variable found")
	}
	
	return &Config{
		Port: port,
		BaseURL: baseURL,
		GithubToken: githubToken,
		DatabaseURL: databaseURL,
	}, nil
}

