package api

import (
	"net/http"
	"runtime"
)

type versionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	GoVersion string `json:"go_version"`
}

// normalize fills unset build values so the response never contains empty
// fields when the binary is built without -ldflags.
func (b BuildInfo) normalize() BuildInfo {
	if b.Version == "" {
		b.Version = "dev"
	}
	if b.Commit == "" {
		b.Commit = "unknown"
	}
	if b.GoVersion == "" {
		b.GoVersion = runtime.Version()
	}
	return b
}

func handleVersion(build BuildInfo) http.HandlerFunc {
	b := build.normalize()
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, versionResponse{
			Version:   b.Version,
			Commit:    b.Commit,
			GoVersion: b.GoVersion,
		})
	}
}
