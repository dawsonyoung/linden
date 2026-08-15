// Package inference abstracts a local model backend. It defines the contract a
// provider must satisfy; no provider is implemented here.
//
// Errors use the taxonomy in src/errs. A provider that cannot be reached is
// unavailable, a provider that does not answer in time is deadline_exceeded,
// and a model the provider does not have is not_found.
package inference

import "context"

// Role identifies who produced a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is one turn in a conversation.
type Message struct {
	Role    Role
	Content string
}

// Model is a model the provider can serve.
type Model struct {
	Name string
}

// Request asks a model to continue a conversation.
type Request struct {
	Model    string
	Messages []Message
}

// Chunk is an incremental fragment of a reply.
type Chunk struct {
	Text string
}

// FinishReason explains why generation ended.
type FinishReason string

const (
	// FinishStop is normal completion.
	FinishStop FinishReason = "stop"
	// FinishCanceled is termination by context cancellation.
	FinishCanceled FinishReason = "canceled"
)

// Result describes a completed generation.
type Result struct {
	Model        string
	FinishReason FinishReason
}

// Client is a model backend.
//
// Implementations must satisfy the conformance suite in
// validation/contracts/inference_contract_test.go.
type Client interface {
	// ListModels reports the models the provider can serve.
	ListModels(ctx context.Context) ([]Model, error)

	// ChatStream generates a reply, delivering it to onChunk as it arrives.
	//
	// onChunk is invoked on the calling goroutine, in order, and only before
	// ChatStream returns. Consumers need no synchronization.
	//
	// If onChunk returns an error, ChatStream stops delivering, performs no
	// further provider work, and returns that error unwrapped so the caller can
	// match its own sentinel.
	//
	// An implementation must check ctx.Err() before each delivery and before
	// classifying any transport failure: cancellation takes precedence over
	// transport classification. Cancellation is user-initiated stop and returns
	// FinishCanceled together with the context error, which is distinct from an
	// onChunk error, which is consumer-side failure.
	//
	// Result is meaningful only when the error is nil or is a context error. On
	// every other error path Result is the zero value.
	//
	// A request carrying no messages is rejected with errs.InvalidArgument
	// before any provider call.
	//
	// Implementations must not log message content at any level.
	//
	// TODO(A.7): a context error carries no errs.Code, so CodeOf reports
	// Internal and a user pressing stop would map to HTTP 500. The taxonomy has
	// no canceled code; resolve when api maps codes to status.
	ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error)
}
