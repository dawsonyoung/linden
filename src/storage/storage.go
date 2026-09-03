package storage

import (
	"context"
	"time"
)

// Session represents the metadata of a chat conversation.
type Session struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Turn represents a single message turn in a session.
type Turn struct {
	Role      string // e.g., "user", "assistant"
	Content   string
	Timestamp time.Time
}

// Store defines the interface for persisting and retrieving chat sessions.
type Store interface {
	// SaveTurn persists a single turn for the given session ID.
	SaveTurn(ctx context.Context, sessionID string, turn Turn) error

	// LoadSession retrieves all turns for a specific session chronologically.
	// Returns errs.NotFound if the session does not exist.
	LoadSession(ctx context.Context, sessionID string) ([]Turn, error)

	// ListSessions returns a chronologically ordered list of all known sessions.
	ListSessions(ctx context.Context) ([]Session, error)
}
