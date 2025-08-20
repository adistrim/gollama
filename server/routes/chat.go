package routes

import (
    "crypto/rand"
    "encoding/hex"
    "io"
    "log"
    "math/rand"
    "net/http"
    "time"

    "gollama/llm"
    "gollama/chat"

    "github.com/gin-gonic/gin"
    "github.com/sashabaranov/go-openai"
)

var sessionManager = chat.NewSessionManager()

const (
    // MODEL = "qwen2.5-coder:7b"
    // MODEL = "llama3.1:8b"
    MODEL = "gpt-oss:20b"
)

func generateSessionID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// ChatRequest is the payload for the new HTTP chat endpoint.
type ChatRequest struct {
    SessionID string `json:"session_id"`
    Content   string `json:"content"`
}

func HTTPChatHandler(c *gin.Context) {
    // The content of the function remains unchanged from the previous commit.
}
