package orchestrator

import (
	"context"
	"errors"
	"testing"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
)

type mockInferenceClient struct {
	ListModelsFunc func(ctx context.Context) ([]inference.Model, error)
	ChatStreamFunc func(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error)
}

func (m *mockInferenceClient) ListModels(ctx context.Context) ([]inference.Model, error) {
	if m.ListModelsFunc != nil {
		return m.ListModelsFunc(ctx)
	}
	return nil, nil
}

func (m *mockInferenceClient) ChatStream(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error) {
	if m.ChatStreamFunc != nil {
		return m.ChatStreamFunc(ctx, req, onChunk)
	}
	return inference.Result{}, nil
}

func Test_Service_ListModels_Success(t *testing.T) {
	client := &mockInferenceClient{
		ListModelsFunc: func(ctx context.Context) ([]inference.Model, error) {
			return []inference.Model{{Name: "test-model"}}, nil
		},
	}
	svc := NewService(client)

	models, err := svc.ListModels(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(models) != 1 || models[0].Name != "test-model" {
		t.Errorf("unexpected models: %+v", models)
	}
}

func Test_Service_ListModels_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("backend failure")
	client := &mockInferenceClient{
		ListModelsFunc: func(ctx context.Context) ([]inference.Model, error) {
			return nil, expectedErr
		},
	}
	svc := NewService(client)

	_, err := svc.ListModels(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func Test_Service_ChatStream_EmptyMessages_ReturnsInvalidArgument(t *testing.T) {
	svc := NewService(&mockInferenceClient{})

	_, err := svc.ChatStream(context.Background(), Request{Model: "test-model"}, func(c Chunk) error { return nil })
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errs.CodeOf(err) != errs.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", errs.CodeOf(err))
	}
}

func Test_Service_ChatStream_Success(t *testing.T) {
	client := &mockInferenceClient{
		ChatStreamFunc: func(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error) {
			if req.Model != "test-model" {
				t.Errorf("expected test-model, got %s", req.Model)
			}
			if len(req.Messages) != 1 || req.Messages[0].Content != "hello" {
				t.Errorf("unexpected messages: %+v", req.Messages)
			}

			_ = onChunk(inference.Chunk{Text: "world"})
			return inference.Result{Model: "test-model", FinishReason: inference.FinishStop}, nil
		},
	}
	svc := NewService(client)

	req := Request{
		Model:    "test-model",
		Messages: []Message{{Role: RoleUser, Content: "hello"}},
	}

	chunks := []string{}
	res, err := svc.ChatStream(context.Background(), req, func(c Chunk) error {
		chunks = append(chunks, c.Text)
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(chunks) != 1 || chunks[0] != "world" {
		t.Errorf("unexpected chunks: %v", chunks)
	}
	if res.Model != "test-model" || res.FinishReason != FinishStop {
		t.Errorf("unexpected result: %+v", res)
	}
}

func Test_Service_ChatStream_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("backend failure")
	client := &mockInferenceClient{
		ChatStreamFunc: func(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error) {
			return inference.Result{Model: "test-model", FinishReason: inference.FinishCanceled}, expectedErr
		},
	}
	svc := NewService(client)

	req := Request{
		Model:    "test-model",
		Messages: []Message{{Role: RoleUser, Content: "hello"}},
	}

	res, err := svc.ChatStream(context.Background(), req, func(c Chunk) error { return nil })

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
	if res.FinishReason != FinishCanceled {
		t.Errorf("expected canceled finish reason, got %v", res.FinishReason)
	}
}
