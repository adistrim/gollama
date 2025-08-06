package chat

import (
	"gollama/llm"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

type ChatSession struct {
	ID       string
	Messages []openai.ChatCompletionMessage
	LastUsed time.Time
	mu       sync.RWMutex
}

func (s *ChatSession) AddMessage(message openai.ChatCompletionMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, message)
	s.LastUsed = time.Now()
}

func (s *ChatSession) GetMessages() []openai.ChatCompletionMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := make([]openai.ChatCompletionMessage, len(s.Messages))
	copy(messages, s.Messages)
	return messages
}

type SessionManager struct {
	sessions map[string]*ChatSession
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*ChatSession),
	}

	go sm.cleanupSessions()
	
	return sm
}

func (sm *SessionManager) GetOrCreateSession(sessionID string) *ChatSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if session, exists := sm.sessions[sessionID]; exists {
		session.LastUsed = time.Now()
		return session
	}
	
	session := &ChatSession{
		ID:       sessionID,
		Messages: make([]openai.ChatCompletionMessage, 0),
		LastUsed: time.Now(),
	}
	
	session.Messages = append(session.Messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: llm.SystemPrompt,
	})
	
	sm.sessions[sessionID] = session
	return session
}

func (sm *SessionManager) cleanupSessions() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for id, session := range sm.sessions {
			// remove sessions inactive for more than 2 hours
			if now.Sub(session.LastUsed) > 2*time.Hour {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}
