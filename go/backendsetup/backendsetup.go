// Package backendsetup configures backend hosts to receive IPIP/IP6IP6
// encapsulated traffic from a katran load balancer.
//
// It handles:
//   - Adding VIP addresses to the loopback interface
//   - Bringing up tunnel interfaces (tunl0 for IPv4, ip6tnl0 for IPv6)
//   - Disabling reverse path filtering on tunnel and specified interfaces
//   - Disabling ARP proxying on tunnel and specified interfaces
package backendsetup

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Config controls which backend setup operations to perform.
type Config struct {
	// VIPs is the list of IP addresses to add to the loopback interface.
	// Both IPv4 and IPv6 addresses are supported.
	VIPs []string
	// Interfaces is the list of interfaces on which to disable rp_filter
	// and proxy_arp. Defaults to ["eth0"] if empty.
	Interfaces []string
	// DryRun prints commands without executing them.
	DryRun bool
}

// Setup performs all backend network configuration needed to receive
// IPIP/IP6IP6 encapsulated traffic from a katran load balancer.
//
// It performs the following steps:
//  1. Ensures tunl0/ip6tnl0 tunnel interfaces exist and are up
//  2. Adds VIP addresses to the loopback interface
//  3. Disables rp_filter on tunnel and specified interfaces
//  4. Disables ARP proxying on tunnel and specified interfaces
//
// Parameters:
//   - cfg: configuration specifying VIPs, interfaces, and dry-run mode.
//
// Returns an error if any operation fails.
func Setup(cfg Config) error {
	ifaces := cfg.Interfaces
	if len(ifaces) == 0 {
		ifaces = []string{"eth0"}
	}

	needV4, needV6 := detectTunnelNeeds(cfg.VIPs)

	// Step 1: bring up tunnel interfaces.
	if needV4 {
		if err := ensureTunnelUp("tunl0", cfg.DryRun); err != nil {
			return fmt.Errorf("bringing up tunl0: %w", err)
		}
	}
	if needV6 {
		if err := ensureTunnelUp("ip6tnl0", cfg.DryRun); err != nil {
			return fmt.Errorf("bringing up ip6tnl0: %w", err)
		}
	}

	// Step 2: add VIP addresses to loopback.
	for _, vip := range cfg.VIPs {
		if err := addVIPToLoopback(vip, cfg.DryRun); err != nil {
			return fmt.Errorf("adding VIP %s to loopback: %w", vip, err)
		}
	}

	// Step 3 & 4: disable rp_filter and proxy_arp on tunnel + specified interfaces.
	var tunIfaces []string
	if needV4 {
		tunIfaces = append(tunIfaces, "tunl0")
	}
	if needV6 {
		tunIfaces = append(tunIfaces, "ip6tnl0")
	}
	allIfaces := append(tunIfaces, ifaces...)

	for _, iface := range allIfaces {
		if err := disableRPFilter(iface, cfg.DryRun); err != nil {
			return fmt.Errorf("disabling rp_filter on %s: %w", iface, err)
		}
		if err := disableProxyARP(iface, cfg.DryRun); err != nil {
			return fmt.Errorf("disabling proxy_arp on %s: %w", iface, err)
		}
	}

	return nil
}

// Teardown removes VIP addresses from the loopback interface.
//
// Parameters:
//   - cfg: configuration specifying VIPs and dry-run mode.
//
// Returns an error if any operation fails.
func Teardown(cfg Config) error {
	for _, vip := range cfg.VIPs {
		if err := removeVIPFromLoopback(vip, cfg.DryRun); err != nil {
			return fmt.Errorf("removing VIP %s from loopback: %w", vip, err)
		}
	}
	return nil
}

// ParseVIPsFromConfig reads a katran YAML config file and extracts
// the unique VIP addresses from the vips section.
//
// Parameters:
//   - path: path to the YAML configuration file.
//
// Returns the deduplicated list of VIP addresses or an error.
func ParseVIPsFromConfig(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	return extractVIPAddresses(data)
}

// DeduplicateVIPs returns a deduplicated copy of the VIP address list,
// preserving the order of first occurrence.
//
// Parameters:
//   - vips: the list of VIP addresses to deduplicate.
//
// Returns a new slice with duplicate addresses removed.
func DeduplicateVIPs(vips []string) []string {
	seen := make(map[string]struct{}, len(vips))
	result := make([]string, 0, len(vips))
	for _, vip := range vips {
		if _, ok := seen[vip]; !ok {
			seen[vip] = struct{}{}
			result = append(result, vip)
		}
	}
	return result
}

// IsIPv6 reports whether the given address string is an IPv6 address.
//
// Parameters:
//   - addr: the IP address string to check.
//
// Returns true if the address is IPv6, false otherwise.
func IsIPv6(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

// detectTunnelNeeds scans the VIP list and returns whether IPv4 and/or
// IPv6 tunnel interfaces are needed.
func detectTunnelNeeds(vips []string) (needV4, needV6 bool) {
	for _, vip := range vips {
		if IsIPv6(vip) {
			needV6 = true
		} else {
			needV4 = true
		}
		if needV4 && needV6 {
			break
		}
	}
	return
}

// ensureTunnelUp brings up the specified tunnel interface.
//
// iface: the tunnel interface name (e.g., "tunl0" or "ip6tnl0").
// dryRun: if true, print the command without executing.
func ensureTunnelUp(iface string, dryRun bool) error {
	return runIP(dryRun, "link", "set", iface, "up")
}

// addVIPToLoopback adds a VIP address to the loopback interface.
// IPv4 addresses get a /32 prefix, IPv6 addresses get /128.
//
// vip: the VIP address to add.
// dryRun: if true, print the command without executing.
func addVIPToLoopback(vip string, dryRun bool) error {
	prefix := "/32"
	if IsIPv6(vip) {
		prefix = "/128"
	}
	return runIP(dryRun, "addr", "add", vip+prefix, "dev", "lo")
}

// removeVIPFromLoopback removes a VIP address from the loopback interface.
//
// vip: the VIP address to remove.
// dryRun: if true, print the command without executing.
func removeVIPFromLoopback(vip string, dryRun bool) error {
	prefix := "/32"
	if IsIPv6(vip) {
		prefix = "/128"
	}
	return runIP(dryRun, "addr", "del", vip+prefix, "dev", "lo")
}

// disableRPFilter writes 0 to /proc/sys/net/ipv4/conf/<iface>/rp_filter.
//
// iface: the interface name.
// dryRun: if true, print the action without writing.
func disableRPFilter(iface string, dryRun bool) error {
	return writeSysctl(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/rp_filter", iface), "0", dryRun)
}

// disableProxyARP writes 0 to /proc/sys/net/ipv4/conf/<iface>/proxy_arp.
//
// iface: the interface name.
// dryRun: if true, print the action without writing.
func disableProxyARP(iface string, dryRun bool) error {
	return writeSysctl(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/proxy_arp", iface), "0", dryRun)
}

// runIP executes an `ip` command with the given arguments.
//
// dryRun: if true, print the command without executing.
// args: the arguments to pass to the `ip` command.
func runIP(dryRun bool, args ...string) error {
	if dryRun {
		fmt.Printf("[dry-run] ip %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.Command("ip", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// writeSysctl writes a value to a sysctl proc file.
//
// path: the full path to the sysctl file.
// value: the value to write.
// dryRun: if true, print the action without writing.
func writeSysctl(path, value string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[dry-run] echo %s > %s\n", value, path)
		return nil
	}
	return os.WriteFile(path, []byte(value), 0644)
}

// BuildIPCommand returns the full command-line arguments that would be
// passed to the `ip` command for a given VIP add/remove operation.
// This is useful for testing command construction.
//
// Parameters:
//   - action: "add" or "del".
//   - vip: the VIP address.
//
// Returns the argument slice for exec.Command("ip", args...).
func BuildIPCommand(action, vip string) []string {
	prefix := "/32"
	if IsIPv6(vip) {
		prefix = "/128"
	}
	return []string{"addr", action, vip + prefix, "dev", "lo"}
}

// BuildTunnelCommand returns the command-line arguments that would be
// passed to the `ip` command to bring up a tunnel interface.
//
// Parameters:
//   - iface: the tunnel interface name.
//
// Returns the argument slice for exec.Command("ip", args...).
func BuildTunnelCommand(iface string) []string {
	return []string{"link", "set", iface, "up"}
}
