package routes

import (
	"log"
	"net/http"

	"gollama/llm"

	"github.com/gin-gonic/gin"
)

type ChatPayload struct {
	Message string `json:"message" binding:"required"`
}

const (
	// MODEL = "qwen2.5-coder:7b"
	MODEL = "llama3.1:8b"
	// MODEL = "gpt-oss:20b"
)

func ChatHandler(c *gin.Context) {
	var payload ChatPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'message' field is required."})
		return
	}

	agent, err := llm.GetAgent(MODEL)
	if err != nil {
		log.Printf("Error creating agent: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize LLM agent."})
		return
	}

	response, err := agent.RunConversation(c.Request.Context(), payload.Message)
	if err != nil {
		log.Printf("Error during conversation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while processing your request."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
