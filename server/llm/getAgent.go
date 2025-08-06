package llm

import (
	"context"
	"fmt"
	"log"
	"sync"

	system "gollama/config"

	"github.com/sashabaranov/go-openai"
)

var (
	instance *Agent
	mu sync.Mutex
)

func GetAgent(modelName string) (*Agent, error) {
	mu.Lock()
	defer mu.Unlock()
	
	if instance != nil {
		if instance.model != modelName {
			instance.model = modelName
		}
		return instance, nil
	}
	
	// token is not needed for local Ollama, but required in openai library
	config := openai.DefaultConfig("") 
	config.BaseURL = system.ENV.BaseURL

	client := openai.NewClientWithConfig(config)

	_, err := client.ListModels(context.Background())
	if err != nil {
		log.Printf("Could not connect to Ollama at %s. Is it running?", config.BaseURL)
		return nil, fmt.Errorf("failed to connect to ollama: %w", err)
	}
	log.Println("Successfully connected to Ollama.")

	instance = &Agent{
		client: client,
		model:  modelName,
	}
	
	return instance, nil
}
