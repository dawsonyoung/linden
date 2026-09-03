//go:build contracts

package contracts

import (
	"context"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/storage"
)

// runStoreContract defines the expected behavior of any storage.Store implementation.
func runStoreContract(t *testing.T, newStore func(t *testing.T) storage.Store) {
	t.Run("A saved turn can be durably retrieved", func(t *testing.T) {
		ctx := context.Background()
		store := newStore(t)

		turn := storage.Turn{
			Role:      "user",
			Content:   "Hello, world",
			Timestamp: time.Now(),
		}

		err := store.SaveTurn(ctx, "session-1", turn)
		if err != nil {
			t.Fatalf("unexpected error saving turn: %v", err)
		}

		turns, err := store.LoadSession(ctx, "session-1")
		if err != nil {
			t.Fatalf("unexpected error loading session: %v", err)
		}

		if len(turns) != 1 {
			t.Fatalf("expected 1 turn, got %d", len(turns))
		}
		if turns[0].Content != "Hello, world" {
			t.Errorf("expected content 'Hello, world', got '%s'", turns[0].Content)
		}
	})

	t.Run("Listing sessions returns chronologically ordered metadata", func(t *testing.T) {
		ctx := context.Background()
		store := newStore(t)

		t1 := time.Now().Add(-2 * time.Hour)
		t2 := time.Now().Add(-1 * time.Hour)

		_ = store.SaveTurn(ctx, "session-old", storage.Turn{Role: "user", Content: "old message", Timestamp: t1})
		_ = store.SaveTurn(ctx, "session-new", storage.Turn{Role: "user", Content: "new message", Timestamp: t2})

		sessions, err := store.ListSessions(ctx)
		if err != nil {
			t.Fatalf("unexpected error listing sessions: %v", err)
		}

		if len(sessions) != 2 {
			t.Fatalf("expected 2 sessions, got %d", len(sessions))
		}
		if sessions[0].ID != "session-old" {
			t.Errorf("expected first session to be 'session-old', got '%s'", sessions[0].ID)
		}
		if sessions[1].ID != "session-new" {
			t.Errorf("expected second session to be 'session-new', got '%s'", sessions[1].ID)
		}
	})

	t.Run("Querying unknown sessions returns errs.NotFound", func(t *testing.T) {
		ctx := context.Background()
		store := newStore(t)

		_, err := store.LoadSession(ctx, "unknown-session")
		if !errs.Is(err, errs.NotFound) {
			t.Errorf("expected NotFound error, got %v", err)
		}
	})
}
