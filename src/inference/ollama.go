package inference

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dawsonyoung/linden/errs"
)

const (
	defaultResponseTimeout = 30 * time.Second
	// A single NDJSON frame. Generous for a text chunk, bounded so a runaway
	// provider cannot exhaust memory.
	maxFrameBytes = 1 << 20
)

// OllamaConfig configures the Ollama adapter.
type OllamaConfig struct {
	// BaseURL is the provider root, for example http://localhost:11434.
	BaseURL string

	// ResponseTimeout bounds the wait for response headers. It does not bound
	// how long a stream may run: generation legitimately takes minutes, and the
	// caller's context governs total duration.
	ResponseTimeout time.Duration
}

// Ollama is a Client backed by an Ollama server.
type Ollama struct {
	baseURL string
	http    *http.Client
}

// NewOllama builds an adapter for the provider at cfg.BaseURL.
func NewOllama(cfg OllamaConfig) (*Ollama, error) {
	if cfg.BaseURL == "" {
		return nil, errs.New(errs.InvalidArgument, "ollama base URL is empty")
	}
	parsed, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, errs.Wrap(errs.InvalidArgument, "parse ollama base URL", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errs.New(errs.InvalidArgument, "ollama base URL scheme must be http or https")
	}

	timeout := cfg.ResponseTimeout
	if timeout <= 0 {
		timeout = defaultResponseTimeout
	}

	return &Ollama{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		http: &http.Client{
			// No Client.Timeout: it would cap the streaming body read.
			Transport: &http.Transport{
				ResponseHeaderTimeout: timeout,
				// Proxy is deliberately nil. Request bodies carry conversation
				// content, and ProxyFromEnvironment would route it to whatever
				// HTTP_PROXY names for any non-loopback provider address.
				Proxy: nil,
			},
		},
	}, nil
}

// ListModels reports the models the provider has installed.
func (o *Ollama) ListModels(ctx context.Context) ([]Model, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "build list models request", err)
	}

	resp, err := o.http.Do(req)
	if err != nil {
		return nil, classifyTransport(ctx, "list models", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// A 404 here means the base URL is not an Ollama server, not that a
		// model is missing.
		if resp.StatusCode == http.StatusNotFound {
			drain(resp)
			return nil, errs.New(errs.Unavailable, "list models: provider endpoint not found, check OLLAMA_URL")
		}
		return nil, classifyStatus("list models", resp)
	}

	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, errs.Wrap(errs.Internal, "decode model list", err)
	}

	models := make([]Model, 0, len(payload.Models))
	for _, m := range payload.Models {
		models = append(models, Model{Name: m.Name})
	}
	return models, nil
}

// ChatStream generates a reply, delivering fragments to onChunk as they arrive.
func (o *Ollama) ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error) {
	if len(req.Messages) == 0 {
		return Result{}, errs.New(errs.InvalidArgument, "request carries no messages")
	}
	if err := ctx.Err(); err != nil {
		return Result{Model: req.Model, FinishReason: FinishCanceled}, err
	}

	body, err := encodeChatRequest(req)
	if err != nil {
		return Result{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return Result{}, errs.Wrap(errs.Internal, "build chat request", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.http.Do(httpReq)
	if err != nil {
		return chatTransportResult(ctx, req, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			drain(resp)
			return Result{}, errs.New(errs.NotFound, "generate reply: provider does not have the requested model")
		}
		return Result{}, classifyStatus("generate reply", resp)
	}

	return o.consumeStream(ctx, req, resp.Body, onChunk)
}

// consumeStream reads Ollama's newline-delimited JSON frames.
func (o *Ollama) consumeStream(ctx context.Context, req Request, body io.Reader, onChunk func(Chunk) error) (Result, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 4096), maxFrameBytes)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var frame chatFrame
		if err := json.Unmarshal(line, &frame); err != nil {
			return Result{}, errs.Wrap(errs.Internal, "decode chat frame", err)
		}
		if frame.Error != "" {
			// The provider reports a failure mid-stream. Its text may echo the
			// request, so it is not carried into the message.
			return Result{}, errs.New(errs.Internal, "provider reported an error mid-stream")
		}

		if frame.Done {
			// A completed stream outranks a late cancellation: the reply is whole.
			// The terminal frame may itself carry the last fragment.
			if frame.Message.Content != "" {
				if err := onChunk(Chunk{Text: frame.Message.Content}); err != nil {
					return Result{}, err
				}
			}
			return Result{Model: req.Model, FinishReason: FinishStop}, nil
		}

		// Checked before every delivery: cancellation outranks anything else.
		if err := ctx.Err(); err != nil {
			return Result{Model: req.Model, FinishReason: FinishCanceled}, err
		}
		if frame.Message.Content == "" {
			continue
		}
		if err := onChunk(Chunk{Text: frame.Message.Content}); err != nil {
			// Unwrapped so the consumer can match its own sentinel.
			return Result{}, err
		}
	}

	if err := scanner.Err(); err != nil {
		// An oversized frame is a provider producing more than we will buffer,
		// not a transport failure; classifying it Unavailable would tell the
		// caller to retry against a healthy provider.
		if errors.Is(err, bufio.ErrTooLong) {
			return Result{}, errs.Wrap(errs.ResourceExhausted, "read chat frame", err)
		}
		return chatTransportResult(ctx, req, err)
	}
	// A provider that closes mid-generation has gone away; that is not a defect
	// in this code, and the caller may reasonably retry.
	return Result{}, errs.New(errs.Unavailable, "provider closed the stream before completion")
}

type chatFrame struct {
	Model   string `json:"model"`
	Done    bool   `json:"done"`
	Error   string `json:"error"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func encodeChatRequest(req Request) ([]byte, error) {
	type wireMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	payload := struct {
		Model    string        `json:"model"`
		Stream   bool          `json:"stream"`
		Messages []wireMessage `json:"messages"`
	}{Model: req.Model, Stream: true}

	for _, m := range req.Messages {
		payload.Messages = append(payload.Messages, wireMessage{Role: string(m.Role), Content: m.Content})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "encode chat request", err)
	}
	return body, nil
}

// chatTransportResult classifies a transport failure, reporting cancellation as
// such so a user pressing stop is never reported as a provider fault.
func chatTransportResult(ctx context.Context, req Request, err error) (Result, error) {
	if ctx.Err() != nil {
		return Result{Model: req.Model, FinishReason: FinishCanceled}, ctx.Err()
	}
	return Result{}, classifyTransport(ctx, "generate reply", err)
}

func classifyTransport(ctx context.Context, op string, err error) error {
	// Cancellation takes precedence over transport classification.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return errs.Wrap(errs.DeadlineExceeded, op, err)
	}
	return errs.Wrap(errs.Unavailable, op, err)
}

// drain discards a response body so the message never carries provider text,
// which may echo the request.
func drain(resp *http.Response) {
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxFrameBytes))
}

func classifyStatus(op string, resp *http.Response) error {
	drain(resp)

	switch {
	case resp.StatusCode == http.StatusBadRequest:
		return errs.New(errs.InvalidArgument, op+": provider rejected the request")
	case resp.StatusCode == http.StatusTooManyRequests:
		return errs.New(errs.ResourceExhausted, op+": provider is rate limiting")
	case resp.StatusCode >= 500:
		return errs.New(errs.Unavailable, fmt.Sprintf("%s: provider returned status %d", op, resp.StatusCode))
	default:
		return errs.New(errs.Internal, fmt.Sprintf("%s: unexpected provider status %d", op, resp.StatusCode))
	}
}
