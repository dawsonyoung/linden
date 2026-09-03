package storage_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/storage"
)

func TestFileStore_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := storage.NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("failed to create FileStore: %v", err)
	}

	ctx := context.Background()

	t.Run("NotFound mapped correctly on load", func(t *testing.T) {
		_, err := store.LoadSession(ctx, "does-not-exist")
		if !errs.Is(err, errs.NotFound) {
			t.Errorf("expected NotFound error, got %v", err)
		}
	})

	t.Run("Malformed JSON on load", func(t *testing.T) {
		badPath := filepath.Join(tmpDir, "bad-session.jsonl")
		err := os.WriteFile(badPath, []byte("{\"ID\":\"bad-session\"}\n{not valid json}\n"), 0644)
		if err != nil {
			t.Fatalf("failed to write bad json: %v", err)
		}

		_, err = store.LoadSession(ctx, "bad-session")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if errs.CodeOf(err) != errs.Internal {
			t.Errorf("expected Internal error, got %v", errs.CodeOf(err))
		}
	})

	t.Run("SaveTurn on non-existent directory", func(t *testing.T) {
		// Create a store pointing to a bad path to force an open error
		badStore, _ := storage.NewFileStore(filepath.Join(tmpDir, "does-not-exist", "deep"))

		// Wait, NewFileStore creates the dir.
		// Let's just remove the dir out from under it.
		os.RemoveAll(filepath.Join(tmpDir, "does-not-exist"))

		err := badStore.SaveTurn(ctx, "session1", storage.Turn{
			Role:      "user",
			Content:   "hello",
			Timestamp: time.Now(),
		})

		if err == nil {
			t.Fatal("expected error saving to missing directory, got nil")
		}
		if errs.CodeOf(err) != errs.Internal {
			t.Errorf("expected Internal error for I/O failure, got %v", errs.CodeOf(err))
		}
	})
}
