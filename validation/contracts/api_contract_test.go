//go:build contracts

package contracts

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/dawsonyoung/linden/api"
)

// The charset a request identifier must fall within to be safe to record.
var safeRequestID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// These tests assert the surface published in docs/product/spec/40-api-surface.md.
// They exercise the handler chain a client reaches, not individual handlers.

func Test_Health_Get_ReturnsOKWithStatusField(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	resp := do(t, h, http.MethodGet, "/health", nil, nil)

	requireStatus(t, resp, http.StatusOK)
	body := requireJSONObject(t, resp)
	if got := requireStringField(t, body, "status"); got != "ok" {
		t.Errorf("status = %q, want %q", got, "ok")
	}
}

func Test_Version_Get_ReturnsBuildIdentification(t *testing.T) {
	h := newHandler(t, api.BuildInfo{Version: "1.2.3", Commit: "abc1234", GoVersion: "go1.22.0"}, nil)

	resp := do(t, h, http.MethodGet, "/version", nil, nil)

	requireStatus(t, resp, http.StatusOK)
	body := requireJSONObject(t, resp)

	for field, want := range map[string]string{
		"version":    "1.2.3",
		"commit":     "abc1234",
		"go_version": "go1.22.0",
	} {
		if got := requireStringField(t, body, field); got != want {
			t.Errorf("%s = %q, want %q", field, got, want)
		}
	}
}

func Test_Version_UnsetBuild_StillReportsAllFields(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	resp := do(t, h, http.MethodGet, "/version", nil, nil)

	requireStatus(t, resp, http.StatusOK)
	body := requireJSONObject(t, resp)

	// The contract is that these fields are always present and non-empty, not
	// what a build without ldflags reports.
	for _, field := range []string{"version", "commit", "go_version"} {
		requireStringField(t, body, field)
	}
}

func Test_AllResponses_CarryRequestID(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/health"},
		{http.MethodGet, "/version"},
		{http.MethodGet, "/does-not-exist"},
		// 405 is produced by the router inside the chain, so it is the case most
		// likely to regress if middleware ordering changes.
		{http.MethodPost, "/health"},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			resp := do(t, h, tc.method, tc.path, nil, nil)
			requireHeader(t, resp, "X-Request-ID")
		})
	}
}

func Test_RequestID_SafeClientValue_Preserved(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)
	const supplied = "client-trace_42"

	resp := do(t, h, http.MethodGet, "/health", nil, map[string]string{"X-Request-ID": supplied})

	if got := requireHeader(t, resp, "X-Request-ID"); got != supplied {
		t.Errorf("X-Request-ID = %q, want the supplied value %q", got, supplied)
	}
}

func Test_RequestID_UnsafeClientValue_Replaced(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	tests := map[string]string{
		"newline":     "trace\ninjected=true",
		"json escape": `trace","forged":"yes`,
		"whitespace":  "trace with spaces",
		"over length": strings.Repeat("a", 300),
	}

	for name, supplied := range tests {
		t.Run(name, func(t *testing.T) {
			resp := do(t, h, http.MethodGet, "/health", nil, map[string]string{"X-Request-ID": supplied})

			got := requireHeader(t, resp, "X-Request-ID")
			if got == supplied {
				t.Fatalf("unsafe client value was echoed: %q", got)
			}
			// Differing is not sufficient: truncating to a safe prefix would also
			// differ while still letting the client control the logged value.
			if !safeRequestID.MatchString(got) {
				t.Errorf("replacement %q is not within the safe charset", got)
			}
		})
	}
}

func Test_UnknownPath_Returns404(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	resp := do(t, h, http.MethodGet, "/no-such-endpoint", nil, nil)

	requireStatus(t, resp, http.StatusNotFound)
}

func Test_KnownPathWrongMethod_Returns405(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, nil)

	for _, path := range []string{"/health", "/version"} {
		t.Run(path, func(t *testing.T) {
			resp := do(t, h, http.MethodPost, path, nil, nil)
			requireStatus(t, resp, http.StatusMethodNotAllowed)
		})
	}
}
