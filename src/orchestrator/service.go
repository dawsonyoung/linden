package orchestrator

import (
	"context"

	"strings"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/storage"
)

// Service is the concrete implementation of ChatService.
type Service struct {
	client inference.Client
	store  storage.Store
}

// NewService creates a new orchestrator Service.
func NewService(client inference.Client, store storage.Store) *Service {
	return &Service{client: client, store: store}
}

// ListModels reports the models available for chat completion.
func (s *Service) ListModels(ctx context.Context) ([]Model, error) {
	infModels, err := s.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}

	var models []Model
	for _, m := range infModels {
		models = append(models, Model{Name: m.Name})
	}
	return models, nil
}

// ChatStream generates a reply, delivering incremental chunks to onChunk.
func (s *Service) ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error) {
	if len(req.Messages) == 0 {
		return Result{}, errs.New(errs.InvalidArgument, "request carries no messages")
	}

	infReq := inference.Request{
		Model: req.Model,
	}

	if req.SessionID != "" {
		history, err := s.store.LoadSession(ctx, req.SessionID)
		if err != nil && !errs.Is(err, errs.NotFound) {
			return Result{}, errs.Wrap(errs.Internal, "failed to load session history", err)
		}

		for _, turn := range history {
			infReq.Messages = append(infReq.Messages, inference.Message{
				Role:    inference.Role(turn.Role),
				Content: turn.Content,
			})
		}
	}

	for _, m := range req.Messages {
		infReq.Messages = append(infReq.Messages, inference.Message{
			Role:    inference.Role(m.Role),
			Content: m.Content,
		})

		if req.SessionID != "" {
			err := s.store.SaveTurn(ctx, req.SessionID, storage.Turn{
				Role:      string(m.Role),
				Content:   m.Content,
				Timestamp: time.Now(),
			})
			if err != nil {
				return Result{}, errs.Wrap(errs.Internal, "failed to save user turn", err)
			}
		}
	}

	var assistantText strings.Builder
	infRes, err := s.client.ChatStream(ctx, infReq, func(c inference.Chunk) error {
		assistantText.WriteString(c.Text)
		return onChunk(Chunk{Text: c.Text})
	})

	if err != nil {
		// As per contract, return zero value on error paths except for context errors.
		// Let's assume context errors return the partial/cancelled Result.
		if errs.CodeOf(err) == errs.DeadlineExceeded || errs.CodeOf(err) == errs.Unavailable || errs.CodeOf(err) == errs.NotFound {
			return Result{}, err
		}
		// If it's a context cancellation, the inference client returns FinishCanceled with the error.
		return Result{
			Model:        infRes.Model,
			FinishReason: FinishReason(infRes.FinishReason),
		}, err
	}

	if req.SessionID != "" && assistantText.Len() > 0 {
		err := s.store.SaveTurn(ctx, req.SessionID, storage.Turn{
			Role:      string(RoleAssistant),
			Content:   assistantText.String(),
			Timestamp: time.Now(),
		})
		if err != nil {
			return Result{
				Model:        infRes.Model,
				FinishReason: FinishReason(infRes.FinishReason),
			}, errs.Wrap(errs.Internal, "failed to save assistant turn", err)
		}
	}

	return Result{
		Model:        infRes.Model,
		FinishReason: FinishReason(infRes.FinishReason),
	}, nil
}
