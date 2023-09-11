package store

import "time"

type PastesStore interface {
	CreatePaste(content, language string, ttl time.Duration) (*Paste, error)
	Paste(id string) (*Paste, error)
	CountPasteView(id string) (int, error)
}

type Paste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
