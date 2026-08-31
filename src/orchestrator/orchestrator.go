// Package orchestrator coordinates workflow-level request routing, context
// assembly, policy decisions, and interaction with inference and storage.
//
// Errors use the taxonomy in src/errs.
package orchestrator

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

// Model is a model available for chat generation.
type Model struct {
	Name string
}

// Request asks for a streamed conversational reply.
type Request struct {
	Model     string
	Messages  []Message
	SessionID string
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

// ChatService coordinates chat interactions.
//
// Implementations must satisfy the conformance suite in
// validation/contracts/orchestrator_contract_test.go.
type ChatService interface {
	// ListModels reports the models available for chat completion.
	ListModels(ctx context.Context) ([]Model, error)

	// ChatStream generates a reply, delivering incremental chunks to onChunk.
	//
	// onChunk is invoked on the calling goroutine, in order, and only before
	// ChatStream returns. Consumers need no synchronization.
	//
	// If onChunk returns an error, ChatStream stops delivering, performs no
	// further work, and returns that error unwrapped so the caller can match its
	// own sentinel.
	//
	// An implementation must check ctx.Err() before each delivery and before
	// classifying any failure: cancellation takes precedence over transport
	// errors. Cancellation returns FinishCanceled and context error.
	//
	// Result is meaningful only when error is nil or context error. On every
	// other error path Result is the zero value.
	//
	// A request carrying no messages is rejected with errs.InvalidArgument
	// before any backend call.
	//
	// Implementations must not log message content at any level.
	ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error)
}
