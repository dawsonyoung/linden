//go:build contracts

package contracts

import (
	"context"
	"testing"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
	"github.com/dawsonyoung/linden/storage"
)

// Test_Service_SatisfiesChatServiceContract runs the conformance suite against
// the concrete Service implementation wired to a scripted inference double.
func Test_Service_SatisfiesChatServiceContract(t *testing.T) {
	runChatServiceContract(t, func(t *testing.T, s orchestratorScript) orchestrator.ChatService {
		store := &scriptedStore{
			sessions: make(map[string]*storage.Session),
			turns:    make(map[string][]storage.Turn),
		}
		return orchestrator.NewService(&scriptedInferenceClient{script: s}, store)
	})
}

// scriptedInferenceClient maps the orchestratorScript behaviors into inference-layer behaviors.
type scriptedInferenceClient struct {
	script orchestratorScript
}

func (c *scriptedInferenceClient) ListModels(ctx context.Context) ([]inference.Model, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := c.providerFailure("list models"); err != nil {
		return nil, err
	}
	var models []inference.Model
	for _, m := range c.script.Models {
		models = append(models, inference.Model{Name: m.Name})
	}
	return models, nil
}

func (c *scriptedInferenceClient) ChatStream(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error) {
	if len(req.Messages) == 0 {
		return inference.Result{}, errs.New(errs.InvalidArgument, "request carries no messages")
	}
	if err := c.providerFailure("generate reply"); err != nil {
		return inference.Result{}, err
	}
	if !c.serves(req.Model) {
		return inference.Result{}, errs.New(errs.NotFound, "requested model is not available")
	}

	for _, text := range c.script.Chunks {
		if err := ctx.Err(); err != nil {
			return inference.Result{Model: req.Model, FinishReason: inference.FinishCanceled}, err
		}
		if err := onChunk(inference.Chunk{Text: text}); err != nil {
			return inference.Result{}, err
		}
	}

	return inference.Result{Model: req.Model, FinishReason: inference.FinishStop}, nil
}

func (c *scriptedInferenceClient) providerFailure(op string) error {
	switch c.script.Fail {
	case failUnreachable:
		return errs.Wrap(errs.Unavailable, op, errScriptedRefused)
	case failTimeout:
		return errs.Wrap(errs.DeadlineExceeded, op, context.DeadlineExceeded)
	default:
		return nil
	}
}

func (c *scriptedInferenceClient) serves(model string) bool {
	for _, m := range c.script.Models {
		if m.Name == model {
			return true
		}
	}
	return false
}
