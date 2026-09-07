package web

import "embed"

// BuildFS contains the static assets of the SvelteKit application.
//
//go:embed build/*
var BuildFS embed.FS
