package discovery

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/dawsonyoung/linden/errs"
	"github.com/grandcat/zeroconf"
)

type Advertiser interface {
	Start(port int) error
	Stop()
}

type zeroConfAdvertiser struct {
	logger *slog.Logger
	server *zeroconf.Server
	id     string
}

func NewAdvertiser(logger *slog.Logger) Advertiser {
	// Generate a 4-character random hex ID
	b := make([]byte, 2)
	rand.Read(b)
	id := hex.EncodeToString(b)

	return &zeroConfAdvertiser{
		logger: logger,
		id:     id,
	}
}

func (a *zeroConfAdvertiser) Start(port int) error {
	instanceName := fmt.Sprintf("Linden AI (%s)", a.id)
	serviceType := "_linden._tcp"
	domain := "local."

	txtRecords := []string{
		"version=0.1.0",
		"api_path=/",
		fmt.Sprintf("id=%s", a.id),
	}

	server, err := zeroconf.Register(instanceName, serviceType, domain, port, txtRecords, nil)
	if err != nil {
		return errs.Wrap(errs.Internal, "failed to start zeroconf mDNS server", err)
	}

	a.server = server
	a.logger.Info("mDNS advertiser started", slog.String("instance", instanceName), slog.String("service", serviceType), slog.Int("port", port))
	return nil
}

func (a *zeroConfAdvertiser) Stop() {
	if a.server != nil {
		a.server.Shutdown()
		a.logger.Info("mDNS advertiser stopped")
	}
}
