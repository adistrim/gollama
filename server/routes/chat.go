package routes

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"gollama/llm"
	"gollama/chat"
	"gollama/socket"

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

func WebSocketHandler(c *gin.Context) {
	conn, err := socket.NewConnection(c)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	go conn.WritePump()

	conn.SendMessage(socket.Message{
		Response: "Hey! How can I assist you today?",
	})

	conn.ReadPump(func(conn *socket.Connection, msg socket.Message) {
		sessionID := msg.SessionID
		if sessionID == "" {
			sessionID = generateSessionID()
		}

		chatSession := sessionManager.GetOrCreateSession(sessionID)

		userMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg.Content,
		}
		chatSession.AddMessage(userMessage)

		agent, err := llm.GetAgent(MODEL)
		if err != nil {
			log.Printf("Error creating agent: %v", err)
			conn.SendMessage(socket.Message{
				Error: "Failed to initialize LLM agent",
			})
			return
		}

		updatedMessages, err := agent.RunSessionConversation(conn.Context(), chatSession.GetMessages())
		if err != nil {
			log.Printf("Error during conversation: %v", err)
			conn.SendMessage(socket.Message{
				Error: "An error occurred while processing your request",
			})
			return
		}

		for _, message := range updatedMessages[len(chatSession.GetMessages()):] {
			chatSession.AddMessage(message)
		}

		if len(updatedMessages) > 0 {
			lastMessage := updatedMessages[len(updatedMessages)-1]
			if lastMessage.Role == openai.ChatMessageRoleAssistant {
				conn.SendMessage(socket.Message{
					Response:  lastMessage.Content,
					SessionID: sessionID,
				})
			}
		}
	})
}
