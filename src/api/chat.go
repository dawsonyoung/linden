package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/orchestrator"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model     string        `json:"model"`
	Messages  []chatMessage `json:"messages"`
	Stream    bool          `json:"stream"`
	SessionID string        `json:"sessionId,omitempty"`
}

type chatChunk struct {
	Text string `json:"text"`
}

func handleChat(chatService orchestrator.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size to 1MB to prevent DOS via massive payloads
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, errs.Wrap(errs.InvalidArgument, "invalid JSON payload or payload too large", err))
			return
		}

		if req.Model == "" {
			writeError(w, errs.New(errs.InvalidArgument, "model name is required"))
			return
		}
		if len(req.Model) > 100 {
			writeError(w, errs.New(errs.InvalidArgument, "model name is too long"))
			return
		}

		if len(req.Messages) == 0 {
			writeError(w, errs.New(errs.InvalidArgument, "messages array is empty"))
			return
		}
		if len(req.Messages) > 1000 {
			writeError(w, errs.New(errs.InvalidArgument, "too many messages"))
			return
		}

		orchReq := orchestrator.Request{
			Model:     req.Model,
			SessionID: req.SessionID,
			Messages:  make([]orchestrator.Message, len(req.Messages)),
		}
		for i, m := range req.Messages {
			if m.Role != "user" && m.Role != "assistant" && m.Role != "system" {
				writeError(w, errs.New(errs.InvalidArgument, "invalid role: must be user, assistant, or system"))
				return
			}
			if len(m.Content) == 0 {
				writeError(w, errs.New(errs.InvalidArgument, "message content cannot be empty"))
				return
			}
			if len(m.Content) > 64000 {
				writeError(w, errs.New(errs.InvalidArgument, "message content too large"))
				return
			}
			orchReq.Messages[i] = orchestrator.Message{
				Role:    orchestrator.Role(m.Role),
				Content: m.Content,
			}
		}

		rc := http.NewResponseController(w)

		// Clear write deadline for SSE streaming
		_ = rc.SetWriteDeadline(time.Time{})

		encoder := json.NewEncoder(w)

		headersWritten := false
		writeEvent := func(event string, data any) error {
			if !headersWritten {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")
				w.WriteHeader(http.StatusOK)
				rc.Flush()
				headersWritten = true
			}
			if _, err := fmt.Fprintf(w, "event: %s\ndata: ", event); err != nil {
				return err
			}
			if err := encoder.Encode(data); err != nil {
				return err
			}
			if _, err := fmt.Fprint(w, "\n"); err != nil {
				return err
			}
			return rc.Flush()
		}

		onChunk := func(c orchestrator.Chunk) error {
			return writeEvent("message", chatChunk{Text: c.Text})
		}

		_, err := chatService.ChatStream(r.Context(), orchReq, onChunk)
		if err != nil {
			if !headersWritten {
				// No headers written yet, meaning it failed before starting the stream (e.g. NotFound)
				writeError(w, err)
				return
			}
			// Headers were already written, stream error event
			code := errs.CodeOf(err)
			codeStr := "internal"
			if code == errs.NotFound {
				codeStr = "not_found"
			}
			_ = writeEvent("error", errorResponse{
				Error: err.Error(),
				Code:  codeStr,
			})
			return
		}

		if !headersWritten {
			// If we get here and no headers were written, it means it completed with no chunks
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			rc.Flush()
			headersWritten = true
		}

		_ = writeEvent("done", struct{}{})
	}
}
