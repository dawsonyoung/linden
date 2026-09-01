//go:build contracts

package contracts

import (
	"net/http"
	"testing"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/orchestrator"
)

func Test_Models_Get_ReturnsModelsArray(t *testing.T) {
	t.Skip("Pending Stage A.7")
	s := orchestratorScript{
		Models: []orchestrator.Model{{Name: "llama2"}, {Name: "mistral"}},
	}
	h := newHandler(t, api.BuildInfo{}, &scriptedChatService{script: s})

	resp := do(t, h, http.MethodGet, "/models", nil, nil)
	requireStatus(t, resp, http.StatusOK)
	body := requireJSONObject(t, resp)

	rawModels, ok := body["models"]
	if !ok {
		t.Fatalf("response has no 'models' field; got %v", fieldNames(body))
	}
	
	models, ok := rawModels.([]any)
	if !ok {
		t.Fatalf("models field = %T, want array", rawModels)
	}

	if len(models) != 2 {
		t.Fatalf("len(models) = %d, want 2", len(models))
	}

	for i, want := range []string{"llama2", "mistral"} {
		obj, ok := models[i].(map[string]any)
		if !ok {
			t.Fatalf("models[%d] = %T, want JSON object", i, models[i])
		}
		if got := requireStringField(t, obj, "name"); got != want {
			t.Errorf("models[%d].name = %q, want %q", i, got, want)
		}
	}
}

func Test_Models_Get_ProviderFailure_Returns503(t *testing.T) {
	t.Skip("Pending Stage A.7")
	s := orchestratorScript{
		Fail: failUnreachable,
	}
	h := newHandler(t, api.BuildInfo{}, &scriptedChatService{script: s})

	resp := do(t, h, http.MethodGet, "/models", nil, nil)
	requireStatus(t, resp, http.StatusServiceUnavailable)
	body := requireJSONObject(t, resp)
	
	if got := requireStringField(t, body, "error"); got == "" {
		t.Error("error field is missing or empty")
	}
	if got := requireStringField(t, body, "code"); got != "unavailable" {
		t.Errorf("code = %q, want unavailable", got)
	}
}
