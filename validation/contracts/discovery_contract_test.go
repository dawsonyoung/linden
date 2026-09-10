//go:build contracts
// +build contracts

package contracts

import (
	"testing"

	"github.com/dawsonyoung/linden/discovery"
)

// mockAdvertiser implements discovery.Advertiser for contract testing
type mockAdvertiser struct {
	started bool
	stopped bool
	port    int
}

func (m *mockAdvertiser) Start(port int) error {
	m.started = true
	m.port = port
	return nil
}

func (m *mockAdvertiser) Stop() {
	m.stopped = true
}

func Test_Discovery_Contract_AdvertiserInterface(t *testing.T) {
	var adv discovery.Advertiser = &mockAdvertiser{}

	err := adv.Start(8080)
	if err != nil {
		t.Fatalf("Start returned unexpected error: %v", err)
	}

	mock, ok := adv.(*mockAdvertiser)
	if !ok || !mock.started {
		t.Error("Expected Start to mutate mock state")
	}
	if mock.port != 8080 {
		t.Errorf("Expected port 8080, got %d", mock.port)
	}

	adv.Stop()
	if !mock.stopped {
		t.Error("Expected Stop to mutate mock state")
	}
}
