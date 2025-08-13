package routes

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	mathRand "math/rand"

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
				Error:     "Failed to initialize LLM agent",
				SessionID: sessionID,
			})
			return
		}

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
				msgText := messages[mathRand.Intn(len(messages))]
				conn.SendMessage(socket.Message{
					Response:     msgText,
					SessionID:    sessionID,
					IsProcessing: true,
				})
			}
		}

		updatedMessages, err := agent.RunSessionConversation(conn.Context(), chatSession.GetMessages(), statusCallback)
		if err != nil {
			log.Printf("Error during conversation: %v", err)
			conn.SendMessage(socket.Message{
				Error:     "An error occurred while processing your request",
				SessionID: sessionID,
			})
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
			conn.SendMessage(socket.Message{
				Response:     lastAssistantMessage,
				SessionID:    sessionID,
				IsProcessing: false,
			})
		} else {
			conn.SendMessage(socket.Message{
				Response:     "I've completed all requested operations. Check the repository for the changes.",
				SessionID:    sessionID,
				IsProcessing: false,
			})
		}
	})
}
