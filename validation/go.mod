// Package validation holds the contract, integration, security, and smoke
// suites. It is a separate module from src so that test-only dependencies never
// enter the server's dependency graph.
module github.com/dawsonyoung/linden/validation

go 1.25.0

require (
	github.com/dawsonyoung/linden v0.0.0
	github.com/grandcat/zeroconf v1.0.0
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/miekg/dns v1.1.27 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
)

replace github.com/dawsonyoung/linden => ../src
