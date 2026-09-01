//go:build contracts

package contracts

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/orchestrator"
)

// orchestratorScript configures a ChatService double to exhibit specific behaviors.
type orchestratorScript struct {
	Models []orchestrator.Model
	Chunks []string
	Fail   failMode
}

// orchestratorServiceFactory builds a ChatService exhibiting the given script.
type orchestratorServiceFactory func(t *testing.T, s orchestratorScript) orchestrator.ChatService

// runChatServiceContract asserts every guarantee orchestrator.ChatService publishes.
func runChatServiceContract(t *testing.T, newService orchestratorServiceFactory) {
	t.Helper()

	t.Run("ListModels returns the available models", func(t *testing.T) {
		want := []orchestrator.Model{{Name: "tinyllama"}, {Name: "llama3"}}
		svc := newService(t, orchestratorScript{Models: want})

		got, err := svc.ListModels(contractContext(t))
		if err != nil {
			t.Fatalf("ListModels() error = %v, want nil", err)
		}
		if len(got) != len(want) {
			t.Fatalf("got %d models, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i].Name != want[i].Name {
				t.Errorf("model[%d] = %q, want %q", i, got[i].Name, want[i].Name)
			}
		}
	})

	t.Run("ListModels maps an unreachable backend to unavailable", func(t *testing.T) {
		svc := newService(t, orchestratorScript{Fail: failUnreachable})

		_, err := svc.ListModels(contractContext(t))
		if err == nil {
			t.Fatal("ListModels() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.Unavailable {
			t.Errorf("code = %q, want %q", got, errs.Unavailable)
		}
	})

	t.Run("ChatStream delivers chunks in order", func(t *testing.T) {
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Chunks: []string{"Hello", ", ", "world"},
		})

		var got []string
		_, err := svc.ChatStream(contractContext(t), validOrchestratorRequest(), func(ch orchestrator.Chunk) error {
			got = append(got, ch.Text)
			return nil
		})
		if err != nil {
			t.Fatalf("ChatStream() error = %v, want nil", err)
		}

		want := []string{"Hello", ", ", "world"}
		if len(got) != len(want) {
			t.Fatalf("received %d chunks, want %d: %q", len(got), len(want), got)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("chunk[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("ChatStream returns the terminal result for the request", func(t *testing.T) {
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Chunks: []string{"a", "b"},
		})

		res, err := svc.ChatStream(contractContext(t), validOrchestratorRequest(), func(orchestrator.Chunk) error { return nil })
		if err != nil {
			t.Fatalf("ChatStream() error = %v, want nil", err)
		}
		if res.FinishReason != orchestrator.FinishStop {
			t.Errorf("FinishReason = %q, want %q", res.FinishReason, orchestrator.FinishStop)
		}
		if res.Model != "tinyllama" {
			t.Errorf("Model = %q, want %q", res.Model, "tinyllama")
		}
	})

	t.Run("ChatStream does not invoke the callback after returning", func(t *testing.T) {
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Chunks: []string{"a", "b", "c"},
		})

		var delivered atomic.Int64
		if _, err := svc.ChatStream(contractContext(t), validOrchestratorRequest(), func(orchestrator.Chunk) error {
			delivered.Add(1)
			return nil
		}); err != nil {
			t.Fatalf("ChatStream() error = %v, want nil", err)
		}

		atReturn := delivered.Load()
		time.Sleep(20 * time.Millisecond)
		if after := delivered.Load(); after != atReturn {
			t.Errorf("callback ran %d more times after ChatStream returned", after-atReturn)
		}
	})

	t.Run("ChatStream stops when the callback returns an error", func(t *testing.T) {
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Chunks: []string{"first", "second", "third"},
		})

		sentinel := errors.New("consumer gave up")
		delivered := 0

		_, err := svc.ChatStream(contractContext(t), validOrchestratorRequest(), func(orchestrator.Chunk) error {
			delivered++
			return sentinel
		})

		if !errors.Is(err, sentinel) {
			t.Fatalf("ChatStream() error = %v, want the callback's error", err)
		}
		if delivered != 1 {
			t.Errorf("callback ran %d times after returning an error, want 1", delivered)
		}
	})

	t.Run("ChatStream honors context cancellation", func(t *testing.T) {
		chunks := []string{"first", "second", "third", "fourth"}
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Chunks: chunks,
		})

		ctx, cancel := context.WithCancel(contractContext(t))
		defer cancel()

		delivered := 0
		res, err := svc.ChatStream(ctx, validOrchestratorRequest(), func(orchestrator.Chunk) error {
			delivered++
			cancel()
			return nil
		})

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("ChatStream() error = %v, want context.Canceled", err)
		}
		if res.FinishReason != orchestrator.FinishCanceled {
			t.Errorf("FinishReason = %q, want %q", res.FinishReason, orchestrator.FinishCanceled)
		}
		if delivered >= len(chunks) {
			t.Errorf("delivered all %d chunks despite cancellation", delivered)
		}
	})

	t.Run("ChatStream maps a provider timeout to deadline_exceeded", func(t *testing.T) {
		svc := newService(t, orchestratorScript{
			Models: []orchestrator.Model{{Name: "tinyllama"}},
			Fail:   failTimeout,
		})

		_, err := svc.ChatStream(contractContext(t), validOrchestratorRequest(), func(orchestrator.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.DeadlineExceeded {
			t.Errorf("code = %q, want %q", got, errs.DeadlineExceeded)
		}
	})

	t.Run("ChatStream maps an unknown model to not_found", func(t *testing.T) {
		svc := newService(t, orchestratorScript{Models: []orchestrator.Model{{Name: "tinyllama"}}})

		req := validOrchestratorRequest()
		req.Model = "not-installed"

		res, err := svc.ChatStream(contractContext(t), req, func(orchestrator.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.NotFound {
			t.Errorf("code = %q, want %q", got, errs.NotFound)
		}
		if res != (orchestrator.Result{}) {
			t.Errorf("Result = %+v on an error path, want the zero value", res)
		}
	})

	t.Run("ChatStream rejects a request with no messages", func(t *testing.T) {
		svc := newService(t, orchestratorScript{Models: []orchestrator.Model{{Name: "tinyllama"}}})

		req := validOrchestratorRequest()
		req.Messages = nil

		_, err := svc.ChatStream(contractContext(t), req, func(orchestrator.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.InvalidArgument {
			t.Errorf("code = %q, want %q", got, errs.InvalidArgument)
		}
	})
}

func validOrchestratorRequest() orchestrator.Request {
	return orchestrator.Request{
		Model:    "tinyllama",
		Messages: []orchestrator.Message{{Role: orchestrator.RoleUser, Content: "hello"}},
	}
}
