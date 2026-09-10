package discovery

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"

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

	var ips []string
	if ifaces, err := net.Interfaces(); err == nil {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok {
					if ip4 := ipNet.IP.To4(); ip4 != nil {
						ips = append(ips, ip4.String())
					}
				}
			}
		}
	}

	server, err := zeroconf.RegisterProxy(instanceName, serviceType, domain, port, "linden", ips, txtRecords, nil)
	if err != nil {
		// Fall back to system hostname registration if proxy registration fails
		server, err = zeroconf.Register(instanceName, serviceType, domain, port, txtRecords, nil)
		if err != nil {
			return errs.Wrap(errs.Internal, "failed to start zeroconf mDNS server", err)
		}
	}

	a.server = server
	a.logger.Info("mDNS advertiser started", slog.String("instance", instanceName), slog.String("service", serviceType), slog.String("host", "linden.local"), slog.Int("port", port))
	return nil
}

func (a *zeroConfAdvertiser) Stop() {
	if a.server != nil {
		a.server.Shutdown()
		a.logger.Info("mDNS advertiser stopped")
	}
}
