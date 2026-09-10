//go:build integration
// +build integration

package integration

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/discovery"
	"github.com/grandcat/zeroconf"
)

func Test_Discovery_Advertiser_AdvertisesToLoopback(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	
	// Start advertiser
	adv := discovery.NewAdvertiser(logger)
	err := adv.Start(8080)
	if err != nil {
		t.Fatalf("Failed to start advertiser: %v", err)
	}
	defer adv.Stop()

	// Fetch all interfaces to ensure we scan loopback explicitly
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get net interfaces: %v", err)
	}
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIfaces(ifaces))
	if err != nil {
		t.Fatalf("Failed to initialize zeroconf resolver: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Browse for the exact custom service type
	err = resolver.Browse(ctx, "_linden._tcp", "local.", entries)
	if err != nil {
		t.Fatalf("Failed to browse: %v", err)
	}

	found := false
	for entry := range entries {
		t.Logf("Discovered service: %s on port %d with TXT %v", entry.Instance, entry.Port, entry.Text)
		// zeroconf sometimes escapes spaces and parens, so use Contains
		if strings.Contains(entry.Instance, "Linden") {
			found = true
			if entry.Port != 8080 {
				t.Errorf("Expected port 8080, got %d", entry.Port)
			}
			hasVersion := false
			hasID := false
			for _, txt := range entry.Text {
				if strings.HasPrefix(txt, "version=") {
					hasVersion = true
				}
				if strings.HasPrefix(txt, "id=") {
					hasID = true
				}
			}
			if !hasVersion {
				t.Errorf("Missing version TXT record, got: %v", entry.Text)
			}
			if !hasID {
				t.Errorf("Missing id TXT record, got: %v", entry.Text)
			}
			break
		}
	}

	if !found {
		t.Error("Did not discover the _linden._tcp service on loopback")
	}
}
