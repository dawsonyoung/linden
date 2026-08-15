//go:build contracts

package contracts

import (
	"context"
	"errors"
	"testing"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
)

// Runs the conformance suite against a scripted double. No provider exists yet;
// this proves the contract is coherent and satisfiable. A.3 adds the same call
// with a factory returning the Ollama adapter.
func Test_ScriptedClient_SatisfiesClientContract(t *testing.T) {
	runClientContract(t, func(_ *testing.T, s script) inference.Client {
		return &scriptedClient{script: s}
	})
}

// scriptedClient is the minimum implementation that satisfies the contract. It
// exists to check the suite, not to be used outside it.
type scriptedClient struct {
	script script
}

func (c *scriptedClient) ListModels(ctx context.Context) ([]inference.Model, error) {
	// Cancellation takes precedence over transport classification.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := c.providerFailure("list models"); err != nil {
		return nil, err
	}
	return c.script.Models, nil
}

func (c *scriptedClient) ChatStream(ctx context.Context, req inference.Request, onChunk func(inference.Chunk) error) (inference.Result, error) {
	if len(req.Messages) == 0 {
		return inference.Result{}, errs.New(errs.InvalidArgument, "request carries no messages")
	}
	if err := c.providerFailure("generate reply"); err != nil {
		return inference.Result{}, err
	}
	if !c.serves(req.Model) {
		// The model name is a caller-supplied string; it stays out of the message.
		return inference.Result{}, errs.New(errs.NotFound, "requested model is not available")
	}

	for _, text := range c.script.Chunks {
		if err := ctx.Err(); err != nil {
			return inference.Result{Model: req.Model, FinishReason: inference.FinishCanceled}, err
		}
		if err := onChunk(inference.Chunk{Text: text}); err != nil {
			// Returned unwrapped so the consumer can match its own sentinel.
			return inference.Result{}, err
		}
	}

	return inference.Result{Model: req.Model, FinishReason: inference.FinishStop}, nil
}

func (c *scriptedClient) providerFailure(op string) error {
	switch c.script.Fail {
	case failUnreachable:
		// A transport error is unclassified until this layer classifies it.
		return errs.Wrap(errs.Unavailable, op, errScriptedRefused)
	case failTimeout:
		return errs.Wrap(errs.DeadlineExceeded, op, context.DeadlineExceeded)
	default:
		return nil
	}
}

func (c *scriptedClient) serves(model string) bool {
	for _, m := range c.script.Models {
		if m.Name == model {
			return true
		}
	}
	return false
}

var errScriptedRefused = errors.New("connection refused")
