// backend-setup configures a backend host to receive IPIP/IP6IP6
// encapsulated traffic from a katran load balancer.
//
// It adds VIP addresses to the loopback interface, brings up tunnel
// interfaces, and disables rp_filter and proxy_arp on the relevant
// network interfaces.
//
// Usage:
//
//	backend-setup -vips "10.0.0.1,2001:db8::1" [-interfaces "eth0"] [-dry-run] [-teardown]
//	backend-setup -config /path/to/config.yaml [-interfaces "eth0"] [-dry-run] [-teardown]
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tehnerd/vatran/go/backendsetup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main logic. It parses flags, collects VIP addresses
// from -vips and/or -config, and calls Setup or Teardown accordingly.
// Returns an error if any operation fails.
func run() error {
	vipsFlag := flag.String("vips", "", "comma-separated list of VIP addresses")
	configFlag := flag.String("config", "", "path to katran YAML config file to extract VIPs from")
	ifacesFlag := flag.String("interfaces", "eth0", "comma-separated list of interfaces for rp_filter/proxy_arp")
	dryRun := flag.Bool("dry-run", false, "print commands without executing")
	teardown := flag.Bool("teardown", false, "remove VIP addresses from loopback (cleanup mode)")
	flag.Parse()

	if *vipsFlag == "" && *configFlag == "" {
		return fmt.Errorf("at least one of -vips or -config is required")
	}

	// Collect VIPs from flags.
	var vips []string
	if *vipsFlag != "" {
		for _, v := range strings.Split(*vipsFlag, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				vips = append(vips, v)
			}
		}
	}

	// Collect VIPs from config file.
	if *configFlag != "" {
		configVIPs, err := backendsetup.ParseVIPsFromConfig(*configFlag)
		if err != nil {
			return fmt.Errorf("parsing config %s: %w", *configFlag, err)
		}
		vips = append(vips, configVIPs...)
	}

	vips = backendsetup.DeduplicateVIPs(vips)
	if len(vips) == 0 {
		return fmt.Errorf("no VIP addresses found")
	}

	// Parse interfaces.
	var ifaces []string
	for _, iface := range strings.Split(*ifacesFlag, ",") {
		iface = strings.TrimSpace(iface)
		if iface != "" {
			ifaces = append(ifaces, iface)
		}
	}

	cfg := backendsetup.Config{
		VIPs:       vips,
		Interfaces: ifaces,
		DryRun:     *dryRun,
	}

	if *teardown {
		return backendsetup.Teardown(cfg)
	}
	return backendsetup.Setup(cfg)
}
