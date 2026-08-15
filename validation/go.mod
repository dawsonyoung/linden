// Package validation holds the contract, integration, security, and smoke
// suites. It is a separate module from src so that test-only dependencies never
// enter the server's dependency graph.
module github.com/dawsonyoung/linden/validation

go 1.22

require github.com/dawsonyoung/linden v0.0.0

replace github.com/dawsonyoung/linden => ../src
