package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dawsonyoung/linden/errs"
	"github.com/dawsonyoung/linden/orchestrator"
)

type openAIReq struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type openAIChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []openAIChoice `json:"choices"`
}

type openAIChoice struct {
	Index        int         `json:"index"`
	Delta        openAIDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type openAIDelta struct {
	Content string `json:"content,omitempty"`
}

func handleOpenAIChat(chatService orchestrator.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

		var req openAIReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, errs.Wrap(errs.InvalidArgument, "invalid JSON payload", err))
			return
		}

		// Simple validation
		if req.Model == "" || len(req.Messages) == 0 {
			writeError(w, errs.New(errs.InvalidArgument, "model and messages are required"))
			return
		}

		orchReq := orchestrator.Request{
			Model:    req.Model,
			Messages: make([]orchestrator.Message, len(req.Messages)),
		}

		for i, m := range req.Messages {
			orchReq.Messages[i] = orchestrator.Message{
				Role:    orchestrator.Role(m.Role),
				Content: m.Content,
			}
		}

		if req.Stream {
			rc := http.NewResponseController(w)
			_ = rc.SetWriteDeadline(time.Time{})

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			rc.Flush()

			encoder := json.NewEncoder(w)
			id := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())

			onChunk := func(c orchestrator.Chunk) error {
				chunk := openAIChunk{
					ID:      id,
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   req.Model,
					Choices: []openAIChoice{
						{
							Index: 0,
							Delta: openAIDelta{Content: c.Text},
						},
					},
				}
				if _, err := fmt.Fprint(w, "data: "); err != nil {
					return err
				}
				if err := encoder.Encode(chunk); err != nil {
					return err
				}
				if _, err := fmt.Fprint(w, "\n"); err != nil {
					return err
				}
				return rc.Flush()
			}

			_, err := chatService.ChatStream(r.Context(), orchReq, onChunk)
			if err != nil {
				// Because headers are written, we just stop the stream.
				// In a full implementation, we might send an error chunk.
				return
			}

			fmt.Fprint(w, "data: [DONE]\n\n")
			rc.Flush()
		} else {
			// Non-streaming fallback
			writeError(w, errs.New(errs.InvalidArgument, "only streaming is supported currently"))
		}
	}
}
