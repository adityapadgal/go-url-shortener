package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// URLData holds the full URL and expiry
type URLData struct {
	OriginalURL string
	CreatedAt   time.Time
	ExpiresAt   time.Time

	AccessCount int
	LastAccess  time.Time
}

// InMemoryStore is a basic map-based store (will replace with DB later)
type InMemoryStore struct {
	data map[string]URLData
	mu   sync.RWMutex
}

// NewStore initializes the store
func NewStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]URLData),
	}
}

// GenerateShortCode generates a short random code
func GenerateShortCode(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// SaveURL saves a shortened URL
func (s *InMemoryStore) SaveURL(code string, url string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[code] = URLData{
		OriginalURL: url,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(ttl),
	}
}

// GetURL returns the original URL by short code
func (s *InMemoryStore) GetURL(code string) (URLData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.data[code]
	if !ok {
		return URLData{}, errors.New("not found")
	}
	if time.Now().After(data.ExpiresAt) {
		return URLData{}, errors.New("expired")
	}
	return data, nil
}

// GetURLWithAnalytics returns the original URL and updates access stats
func (s *InMemoryStore) GetURLWithAnalytics(code string) (URLData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.data[code]
	if !ok {
		return URLData{}, errors.New("not found")
	}
	if time.Now().After(data.ExpiresAt) {
		return URLData{}, errors.New("expired")
	}

	// Update analytics
	data.AccessCount++
	data.LastAccess = time.Now()
	s.data[code] = data

	return data, nil
}

