package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Chat_Post_InvalidJSON_Returns400(t *testing.T) {
	srv := NewServer("localhost:0", discardLogger(), BuildInfo{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Result().StatusCode, http.StatusBadRequest)
	}
}
