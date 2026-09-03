package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/dawsonyoung/linden/errs"
)

// FileStore implements Store using JSON Lines files on disk.
type FileStore struct {
	baseDir string
	mu      sync.Mutex
}

// NewFileStore creates a FileStore that saves session files to baseDir.
// It creates the directory if it does not exist.
func NewFileStore(baseDir string) (*FileStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, errs.Wrap(errs.Internal, "initialize storage directory", err)
	}
	return &FileStore{baseDir: baseDir}, nil
}

func (s *FileStore) getPath(sessionID string) string {
	// Using a flat file structure. Assume sessionID is safe.
	return filepath.Join(s.baseDir, sessionID+".jsonl")
}

func (s *FileStore) SaveTurn(ctx context.Context, sessionID string, turn Turn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.getPath(sessionID)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errs.Wrap(errs.Internal, "open session file", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return errs.Wrap(errs.Internal, "stat session file", err)
	}

	if info.Size() == 0 {
		sess := Session{
			ID:        sessionID,
			CreatedAt: turn.Timestamp,
			UpdatedAt: turn.Timestamp,
		}
		if err := json.NewEncoder(file).Encode(sess); err != nil {
			return errs.Wrap(errs.Internal, "write session metadata", err)
		}
	}

	if err := json.NewEncoder(file).Encode(turn); err != nil {
		return errs.Wrap(errs.Internal, "write turn", err)
	}

	return nil
}

func (s *FileStore) LoadSession(ctx context.Context, sessionID string) ([]Turn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.getPath(sessionID)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errs.Wrap(errs.NotFound, "load session", err)
		}
		return nil, errs.Wrap(errs.Internal, "open session file", err)
	}
	defer file.Close()

	var turns []Turn
	scanner := bufio.NewScanner(file)
	
	// Skip the first line (Session metadata)
	if !scanner.Scan() {
		return nil, nil
	}

	for scanner.Scan() {
		var turn Turn
		if err := json.Unmarshal(scanner.Bytes(), &turn); err != nil {
			return nil, errs.Wrap(errs.Internal, "parse turn", err)
		}
		turns = append(turns, turn)
	}

	if err := scanner.Err(); err != nil {
		return nil, errs.Wrap(errs.Internal, "read session file", err)
	}

	return turns, nil
}

func (s *FileStore) ListSessions(ctx context.Context) ([]Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "read storage directory", err)
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".jsonl" {
			continue
		}

		path := filepath.Join(s.baseDir, entry.Name())
		file, err := os.Open(path)
		if err != nil {
			continue // skip unreadable files
		}

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			var sess Session
			if err := json.Unmarshal(scanner.Bytes(), &sess); err == nil {
				info, err := entry.Info()
				if err == nil {
					sess.UpdatedAt = info.ModTime()
				}
				sessions = append(sessions, sess)
			}
		}
		file.Close()
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.Before(sessions[j].CreatedAt)
	})

	return sessions, nil
}
