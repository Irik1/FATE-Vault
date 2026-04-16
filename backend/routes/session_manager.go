package routes

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Session struct {
	ID             string
	UserID         string
	CreatedAt      time.Time
	LastActivityAt time.Time
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *InMemorySessionStore) read(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[id]
}

func (s *InMemorySessionStore) write(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
}

func (s *InMemorySessionStore) destroy(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}

func (s *InMemorySessionStore) gc(idleExpiration, absoluteExpiration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		if time.Since(session.LastActivityAt) > idleExpiration ||
			time.Since(session.CreatedAt) > absoluteExpiration {
			delete(s.sessions, id)
		}
	}
}

type SessionManager struct {
	store              *InMemorySessionStore
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func envIntSeconds(key string, fallbackSeconds int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return time.Duration(fallbackSeconds) * time.Second
}

func NewSessionManager(
	store *InMemorySessionStore,
	gcInterval time.Duration,
	idleExpiration time.Duration,
	absoluteExpiration time.Duration,
) *SessionManager {
	m := &SessionManager{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
	}

	go m.gc(gcInterval)

	return m
}

func (m *SessionManager) gc(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		m.store.gc(m.idleExpiration, m.absoluteExpiration)
	}
}

func generateSessionID() string {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic("failed to generate session id")
	}
	return base64.RawURLEncoding.EncodeToString(id)
}

func (m *SessionManager) validate(session *Session) bool {
	if session == nil {
		return false
	}
	if time.Since(session.CreatedAt) > m.absoluteExpiration ||
		time.Since(session.LastActivityAt) > m.idleExpiration {
		m.store.destroy(session.ID)
		return false
	}
	return true
}

func (m *SessionManager) Create(userID string) *Session {
	now := time.Now()
	session := &Session{
		ID:             generateSessionID(),
		UserID:         userID,
		CreatedAt:      now,
		LastActivityAt: now,
	}
	m.store.write(session)
	return session
}

func (m *SessionManager) ReadValid(id string) *Session {
	session := m.store.read(id)
	if !m.validate(session) {
		return nil
	}
	return session
}

func (m *SessionManager) Touch(session *Session) {
	if session == nil {
		return
	}
	session.LastActivityAt = time.Now()
	m.store.write(session)
}

func (m *SessionManager) Destroy(id string) {
	m.store.destroy(id)
}

func (m *SessionManager) IdleExpiration() time.Duration {
	return m.idleExpiration
}

func (m *SessionManager) AbsoluteExpiration() time.Duration {
	return m.absoluteExpiration
}

var sessionManager = NewSessionManager(
	NewInMemorySessionStore(),
	envDuration("SESSION_GC_INTERVAL", 30*time.Minute),
	envDuration("SESSION_IDLE_EXPIRATION", 1*time.Hour),
	envDuration("SESSION_ABSOLUTE_EXPIRATION", 12*time.Hour),
)
