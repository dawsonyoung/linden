package orchestrator

import (
	"context"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
)

// Service is the concrete implementation of ChatService.
type Service struct {
	client inference.Client
}

// NewService creates a new orchestrator Service.
func NewService(client inference.Client) *Service {
	return &Service{client: client}
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
	for _, m := range req.Messages {
		infReq.Messages = append(infReq.Messages, inference.Message{
			Role:    inference.Role(m.Role),
			Content: m.Content,
		})
	}

	infRes, err := s.client.ChatStream(ctx, infReq, func(c inference.Chunk) error {
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

	return Result{
		Model:        infRes.Model,
		FinishReason: FinishReason(infRes.FinishReason),
	}, nil
}
