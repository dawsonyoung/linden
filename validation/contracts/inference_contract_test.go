//go:build contracts

package contracts

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/inference"
)

// script configures an implementation to exhibit a specific behavior, so the
// same suite can run against a scripted double here and against the Ollama
// adapter backed by httptest in A.3.
type script struct {
	// Models is what the provider serves. It drives ListModels and decides
	// which requests are not_found: a model outside this set is unknown.
	Models []inference.Model
	// Chunks is the reply delivered one element at a time.
	Chunks []string
	// Fail makes the provider misbehave in a specific way.
	Fail failMode
}

// contractTimeout bounds every suite call. A provider that stalls must fail the
// subtest rather than hang until the package-level panic.
const contractTimeout = 5 * time.Second

func contractContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), contractTimeout)
	t.Cleanup(cancel)
	return ctx
}

type failMode int

const (
	failNone failMode = iota
	// failUnreachable is a provider that cannot be contacted.
	failUnreachable
	// failTimeout is a provider that does not answer in time.
	failTimeout
)

// clientFactory builds a Client exhibiting the given script.
type clientFactory func(t *testing.T, s script) inference.Client

// runClientContract asserts every guarantee inference.Client publishes.
//
// Any implementation must pass it unchanged. A.3 calls this with a factory that
// returns the Ollama adapter pointed at an httptest server.
func runClientContract(t *testing.T, newClient clientFactory) {
	t.Helper()

	t.Run("ListModels returns the provider's models", func(t *testing.T) {
		want := []inference.Model{{Name: "tinyllama"}, {Name: "llama3"}}
		c := newClient(t, script{Models: want})

		got, err := c.ListModels(contractContext(t))
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

	t.Run("ListModels maps an unreachable provider to unavailable", func(t *testing.T) {
		c := newClient(t, script{Fail: failUnreachable})

		_, err := c.ListModels(contractContext(t))
		if err == nil {
			t.Fatal("ListModels() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.Unavailable {
			t.Errorf("code = %q, want %q", got, errs.Unavailable)
		}
	})

	t.Run("ChatStream delivers chunks in order", func(t *testing.T) {
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Chunks: []string{"Hello", ", ", "world"},
		})

		var got []string
		_, err := c.ChatStream(contractContext(t), validRequest(), func(ch inference.Chunk) error {
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
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Chunks: []string{"a", "b"},
		})

		res, err := c.ChatStream(contractContext(t), validRequest(), func(inference.Chunk) error { return nil })
		if err != nil {
			t.Fatalf("ChatStream() error = %v, want nil", err)
		}
		if res.FinishReason != inference.FinishStop {
			t.Errorf("FinishReason = %q, want %q", res.FinishReason, inference.FinishStop)
		}
		if res.Model != "tinyllama" {
			t.Errorf("Model = %q, want %q", res.Model, "tinyllama")
		}
	})

	t.Run("ChatStream does not invoke the callback after returning", func(t *testing.T) {
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Chunks: []string{"a", "b", "c"},
		})

		// Atomic because a non-conforming implementation would deliver from
		// another goroutine, and the test must report that rather than race.
		var delivered atomic.Int64
		if _, err := c.ChatStream(contractContext(t), validRequest(), func(inference.Chunk) error {
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
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Chunks: []string{"first", "second", "third"},
		})

		sentinel := errors.New("consumer gave up")
		delivered := 0

		_, err := c.ChatStream(contractContext(t), validRequest(), func(inference.Chunk) error {
			delivered++
			return sentinel
		})

		// Unwrapped, so a consumer can match its own sentinel.
		if !errors.Is(err, sentinel) {
			t.Fatalf("ChatStream() error = %v, want the callback's error", err)
		}
		if delivered != 1 {
			t.Errorf("callback ran %d times after returning an error, want 1", delivered)
		}
	})

	t.Run("ChatStream honors context cancellation", func(t *testing.T) {
		chunks := []string{"first", "second", "third", "fourth"}
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Chunks: chunks,
		})

		ctx, cancel := context.WithCancel(contractContext(t))
		defer cancel()

		delivered := 0
		res, err := c.ChatStream(ctx, validRequest(), func(inference.Chunk) error {
			delivered++
			cancel()
			return nil
		})

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("ChatStream() error = %v, want context.Canceled", err)
		}
		if res.FinishReason != inference.FinishCanceled {
			t.Errorf("FinishReason = %q, want %q", res.FinishReason, inference.FinishCanceled)
		}
		// Not an exact count: an implementation reading a buffered stream may have
		// decoded more than one frame before observing the cancellation. The
		// guarantee is that it stopped early.
		if delivered >= len(chunks) {
			t.Errorf("delivered all %d chunks despite cancellation", delivered)
		}
	})

	t.Run("ChatStream maps a provider timeout to deadline_exceeded", func(t *testing.T) {
		c := newClient(t, script{
			Models: []inference.Model{{Name: "tinyllama"}},
			Fail:   failTimeout,
		})

		_, err := c.ChatStream(contractContext(t), validRequest(), func(inference.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.DeadlineExceeded {
			t.Errorf("code = %q, want %q", got, errs.DeadlineExceeded)
		}
	})

	t.Run("ChatStream maps an unknown model to not_found", func(t *testing.T) {
		c := newClient(t, script{Models: []inference.Model{{Name: "tinyllama"}}})

		req := validRequest()
		req.Model = "not-installed"

		res, err := c.ChatStream(contractContext(t), req, func(inference.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.NotFound {
			t.Errorf("code = %q, want %q", got, errs.NotFound)
		}
		// Result is meaningful only when the error is nil or a context error.
		if res != (inference.Result{}) {
			t.Errorf("Result = %+v on an error path, want the zero value", res)
		}
	})

	t.Run("ChatStream rejects a request with no messages", func(t *testing.T) {
		c := newClient(t, script{Models: []inference.Model{{Name: "tinyllama"}}})

		req := validRequest()
		req.Messages = nil

		_, err := c.ChatStream(contractContext(t), req, func(inference.Chunk) error { return nil })
		if err == nil {
			t.Fatal("ChatStream() error = nil, want an error")
		}
		if got := errs.CodeOf(err); got != errs.InvalidArgument {
			t.Errorf("code = %q, want %q", got, errs.InvalidArgument)
		}
	})
}

func validRequest() inference.Request {
	return inference.Request{
		Model:    "tinyllama",
		Messages: []inference.Message{{Role: inference.RoleUser, Content: "hello"}},
	}
}
