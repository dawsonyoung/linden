//go:build contracts

package contracts

import (
	"context"
	"testing"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/orchestrator"
)

// Runs the conformance suite against a scripted double. This proves the contract is
// coherent and satisfiable in isolation with minimal dependencies.
func Test_ScriptedChatService_SatisfiesChatServiceContract(t *testing.T) {
	runChatServiceContract(t, func(_ *testing.T, s orchestratorScript) orchestrator.ChatService {
		return &scriptedChatService{script: s}
	})
}

// scriptedChatService is the minimum implementation that satisfies the ChatService contract.
type scriptedChatService struct {
	script orchestratorScript
}

func (c *scriptedChatService) ListModels(ctx context.Context) ([]orchestrator.Model, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := c.providerFailure("list models"); err != nil {
		return nil, err
	}
	return c.script.Models, nil
}

func (c *scriptedChatService) ChatStream(ctx context.Context, req orchestrator.Request, onChunk func(orchestrator.Chunk) error) (orchestrator.Result, error) {
	if len(req.Messages) == 0 {
		return orchestrator.Result{}, errs.New(errs.InvalidArgument, "request carries no messages")
	}
	if err := c.providerFailure("generate reply"); err != nil {
		return orchestrator.Result{}, err
	}
	if !c.serves(req.Model) {
		return orchestrator.Result{}, errs.New(errs.NotFound, "requested model is not available")
	}

	for _, text := range c.script.Chunks {
		if err := ctx.Err(); err != nil {
			return orchestrator.Result{Model: req.Model, FinishReason: orchestrator.FinishCanceled}, err
		}
		if err := onChunk(orchestrator.Chunk{Text: text}); err != nil {
			return orchestrator.Result{}, err
		}
	}

	return orchestrator.Result{Model: req.Model, FinishReason: orchestrator.FinishStop}, nil
}

func (c *scriptedChatService) providerFailure(op string) error {
	switch c.script.Fail {
	case failUnreachable:
		return errs.Wrap(errs.Unavailable, op, errScriptedRefused)
	case failTimeout:
		return errs.Wrap(errs.DeadlineExceeded, op, context.DeadlineExceeded)
	default:
		return nil
	}
}

func (c *scriptedChatService) serves(model string) bool {
	for _, m := range c.script.Models {
		if m.Name == model {
			return true
		}
	}
	return false
}
