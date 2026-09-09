package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]time.Time),
	}
}

func (sm *SessionManager) CreateSession() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[token] = time.Now().Add(24 * time.Hour)

	return token, nil
}

func (sm *SessionManager) ValidateSession(token string) bool {
	if token == "" {
		return false
	}

	sm.mu.RLock()
	expiry, exists := sm.sessions[token]
	sm.mu.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		sm.mu.Lock()
		delete(sm.sessions, token)
		sm.mu.Unlock()
		return false
	}

	return true
}

func (sm *SessionManager) InvalidateSession(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, token)
}
