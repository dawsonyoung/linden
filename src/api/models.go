package api

import (
	"encoding/json"
	"net/http"

	"github.com/dawsonyoung/linden/orchestrator"
)

type modelResponse struct {
	Name string `json:"name"`
}

type modelsResponse struct {
	Models []modelResponse `json:"models"`
}

func handleModels(chatService orchestrator.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models, err := chatService.ListModels(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}

		resp := modelsResponse{Models: make([]modelResponse, 0, len(models))}
		for _, m := range models {
			resp.Models = append(resp.Models, modelResponse{Name: m.Name})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
