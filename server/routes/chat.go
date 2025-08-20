package routes

import (
    "crypto/rand"   // for session ID generation
    "encoding/hex"
    "io"
    "log"
    "math/rand"    // local random for message selection
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
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
        return
    }

    sessionID := req.SessionID
    if sessionID == "" {
        sessionID = generateSessionID()
    }

    chatSession := sessionManager.GetOrCreateSession(sessionID)
    userMessage := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: req.Content}
    chatSession.AddMessage(userMessage)

    agent, err := llm.GetAgent(MODEL)
    if err != nil {
        log.Printf("Error creating agent: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize LLM agent"})
        return
    }

    // Set headers for SSE
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.WriteHeader(http.StatusOK)

    // Helper to send SSE messages
    sendEvent := func(event, data string) error {
        if _, err := c.Writer.Write([]byte("event:" + event + "\ndata:" + data + "\n\n")); err != nil {
            return err
        }
        if flusher, ok := c.Writer.(http.Flusher); ok {
            flusher.Flush()
        }
        return nil
    }

    // Initial ack to confirm session start
    sendEvent("info", "Chat session started")

    toolsUsed := false
    statusCallback := func(status string) {
        if !toolsUsed && status == "tools_used" {
            toolsUsed = true
            messages := []string{
                "Working on it… apparently this takes more than two seconds.",
                "Processing… you'll know when I finally survive this step.",
                "Processing… I'll update you shortly.",
                "Working through the steps… hang tight.",
                "The tools and I are having a deep conversation.",
            }
            msgText := messages[mathrand.Intn(len(messages))]
            sendEvent("update", msgText)
        }
    }

    updatedMessages, err := agent.RunSessionConversation(c.Request.Context(), chatSession.GetMessages(), statusCallback)
    if err != nil {
        log.Printf("Error during conversation: %v", err)
        sendEvent("error", "An error occurred while processing your request")
        return
    }

    for _, message := range updatedMessages[len(chatSession.GetMessages()):] {
        chatSession.AddMessage(message)
    }

    var lastAssistantMessage string
    for i := len(updatedMessages) - 1; i >= 0; i-- {
        if updatedMessages[i].Role == openai.ChatMessageRoleAssistant {
            lastAssistantMessage = updatedMessages[i].Content
            break
        }
    }

    if lastAssistantMessage != "" {
        sendEvent("final", lastAssistantMessage)
    } else {
        sendEvent("final", "I've completed all requested operations. Check the repository for the changes.")
    }

    // Slight pause before closing the stream
    time.Sleep(500 * time.Millisecond)
}
