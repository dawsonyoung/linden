// Package validation holds the contract, integration, security, and smoke
// suites. It is a separate module from src so that test-only dependencies never
// enter the server's dependency graph.
module github.com/dawsonyoung/linden/validation

go 1.22

require (
	github.com/dawsonyoung/linden v0.0.0
	github.com/grandcat/zeroconf v1.0.0
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/miekg/dns v1.1.27 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe // indirect
)

replace github.com/dawsonyoung/linden => ../src
