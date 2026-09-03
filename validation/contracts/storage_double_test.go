//go:build contracts

package contracts

import (
	"context"
	"sort"
	"testing"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/storage"
)

func Test_ScriptedStore_SatisfiesStoreContract(t *testing.T) {
	runStoreContract(t, func(t *testing.T) storage.Store {
		return &scriptedStore{
			sessions: make(map[string]*storage.Session),
			turns:    make(map[string][]storage.Turn),
		}
	})
}

// scriptedStore is an in-memory test double for the storage.Store interface.
// It is used to verify the storage contracts before a real persistence layer is built.
type scriptedStore struct {
	sessions map[string]*storage.Session
	turns    map[string][]storage.Turn
}

func (s *scriptedStore) SaveTurn(ctx context.Context, sessionID string, turn storage.Turn) error {
	if _, ok := s.sessions[sessionID]; !ok {
		s.sessions[sessionID] = &storage.Session{
			ID:        sessionID,
			CreatedAt: turn.Timestamp,
			UpdatedAt: turn.Timestamp,
		}
	} else {
		s.sessions[sessionID].UpdatedAt = turn.Timestamp
	}
	s.turns[sessionID] = append(s.turns[sessionID], turn)
	return nil
}

func (s *scriptedStore) LoadSession(ctx context.Context, sessionID string) ([]storage.Turn, error) {
	turns, ok := s.turns[sessionID]
	if !ok {
		return nil, errs.New(errs.NotFound, "session not found")
	}
	return turns, nil
}

func (s *scriptedStore) ListSessions(ctx context.Context) ([]storage.Session, error) {
	var list []storage.Session
	for _, sess := range s.sessions {
		list = append(list, *sess)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}
