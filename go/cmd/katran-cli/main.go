// katran-cli is a command-line interface for the Katran Load Balancer REST API.
//
// Usage:
//
//	katran-cli [global flags] <command> [command flags]
//
// Commands:
//
//	vip     - Manage VIPs (Virtual IPs)
//	real    - Manage real servers (backends)
//	stats   - View statistics
//	mac     - Manage MAC address
//	hc      - Manage healthcheck destinations
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tehnerd/vatran/go/client"
)

var (
	serverURL string
	timeout   int
)

func main() {
	flag.StringVar(&serverURL, "server", "http://localhost:8080", "Katran server URL")
	flag.IntVar(&timeout, "timeout", 30, "Request timeout in seconds")
	flag.Usage = printUsage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	c := client.New(serverURL, client.WithTimeout(time.Duration(timeout)*time.Second))

	cmd := args[0]
	cmdArgs := args[1:]

	var err error
	switch cmd {
	case "vip":
		err = handleVIP(c, cmdArgs)
	case "real":
		err = handleReal(c, cmdArgs)
	case "stats":
		err = handleStats(c, cmdArgs)
	case "mac":
		err = handleMAC(c, cmdArgs)
	case "hc":
		err = handleHC(c, cmdArgs)
	case "health":
		err = handleHealth(c)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Katran Load Balancer CLI

Usage:
  katran-cli [flags] <command> [args]

Global Flags:
  -server string   Katran server URL (default "http://localhost:8080")
  -timeout int     Request timeout in seconds (default 30)

Commands:
  vip     Manage VIPs (list, add, remove, show)
  real    Manage real servers (add, remove, update)
  stats   View statistics (vip, lru, xdp, decap, etc.)
  mac     Manage default router MAC address (show, set)
  hc      Manage healthcheck destinations (list, add, remove)
  health  Check server health

Examples:
  katran-cli vip list
  katran-cli vip add 10.0.0.1 80 tcp
  katran-cli vip show 10.0.0.1 80 tcp
  katran-cli real add 10.0.0.1 80 tcp 192.168.1.1 100
  katran-cli stats vip 10.0.0.1 80 tcp --watch
  katran-cli mac show
  katran-cli hc list
`)
}

// parseProto converts protocol name or number to uint8.
func parseProto(s string) (uint8, error) {
	switch strings.ToLower(s) {
	case "tcp", "6":
		return 6, nil
	case "udp", "17":
		return 17, nil
	default:
		n, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("invalid protocol: %s (use tcp, udp, or number)", s)
		}
		return uint8(n), nil
	}
}

// protoToString converts protocol number to name.
func protoToString(proto uint8) string {
	switch proto {
	case 6:
		return "tcp"
	case 17:
		return "udp"
	default:
		return strconv.Itoa(int(proto))
	}
}

// handleVIP handles VIP management commands.
func handleVIP(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vip <list|add|remove|show> [args]")
	}

	switch args[0] {
	case "list":
		return vipList(c)
	case "add":
		return vipAdd(c, args[1:])
	case "remove", "rm", "delete", "del":
		return vipRemove(c, args[1:])
	case "show", "get":
		return vipShow(c, args[1:])
	default:
		return fmt.Errorf("unknown vip command: %s", args[0])
	}
}

func vipList(c *client.Client) error {
	vips, err := c.ListVIPs()
	if err != nil {
		return err
	}

	if len(vips) == 0 {
		fmt.Println("No VIPs configured")
		return nil
	}

	fmt.Printf("%-40s %-8s %-8s\n", "ADDRESS", "PORT", "PROTO")
	fmt.Println(strings.Repeat("-", 60))
	for _, vip := range vips {
		fmt.Printf("%-40s %-8d %-8s\n", vip.Address, vip.Port, protoToString(vip.Proto))
	}
	return nil
}

func vipAdd(c *client.Client, args []string) error {
	fs := flag.NewFlagSet("vip add", flag.ExitOnError)
	flags := fs.Uint("flags", 0, "VIP flags")
	fs.Parse(args)

	posArgs := fs.Args()
	if len(posArgs) < 3 {
		return fmt.Errorf("usage: vip add <address> <port> <proto> [--flags <flags>]")
	}

	port, err := strconv.ParseUint(posArgs[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	proto, err := parseProto(posArgs[2])
	if err != nil {
		return err
	}

	vip := client.VIP{
		Address: posArgs[0],
		Port:    uint16(port),
		Proto:   proto,
		Flags:   uint32(*flags),
	}

	if err := c.AddVIP(vip); err != nil {
		return err
	}

	fmt.Printf("VIP %s:%d/%s added\n", vip.Address, vip.Port, protoToString(vip.Proto))
	return nil
}

func vipRemove(c *client.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: vip remove <address> <port> <proto>")
	}

	port, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	proto, err := parseProto(args[2])
	if err != nil {
		return err
	}

	if err := c.DeleteVIP(args[0], uint16(port), proto); err != nil {
		return err
	}

	fmt.Printf("VIP %s:%d/%s removed\n", args[0], port, protoToString(proto))
	return nil
}

func vipShow(c *client.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: vip show <address> <port> <proto>")
	}

	port, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	proto, err := parseProto(args[2])
	if err != nil {
		return err
	}

	reals, err := c.GetVIPReals(args[0], uint16(port), proto)
	if err != nil {
		return err
	}

	fmt.Printf("VIP: %s:%d/%s\n", args[0], port, protoToString(proto))
	fmt.Println()

	if len(reals) == 0 {
		fmt.Println("No real servers configured")
		return nil
	}

	fmt.Printf("Real Servers (%d):\n", len(reals))
	fmt.Printf("  %-40s %-10s %-8s\n", "ADDRESS", "WEIGHT", "FLAGS")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for _, real := range reals {
		fmt.Printf("  %-40s %-10d %-8d\n", real.Address, real.Weight, real.Flags)
	}
	return nil
}

// handleReal handles real server management commands.
func handleReal(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: real <add|remove|update> [args]")
	}

	switch args[0] {
	case "add":
		return realAdd(c, args[1:])
	case "remove", "rm", "delete", "del":
		return realRemove(c, args[1:])
	case "update":
		return realUpdate(c, args[1:])
	default:
		return fmt.Errorf("unknown real command: %s", args[0])
	}
}

func realAdd(c *client.Client, args []string) error {
	fs := flag.NewFlagSet("real add", flag.ExitOnError)
	flags := fs.Uint("flags", 0, "Real server flags")
	fs.Parse(args)

	posArgs := fs.Args()
	if len(posArgs) < 5 {
		return fmt.Errorf("usage: real add <vip-addr> <vip-port> <vip-proto> <real-addr> <weight> [--flags <flags>]")
	}

	vipPort, err := strconv.ParseUint(posArgs[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid VIP port: %v", err)
	}

	vipProto, err := parseProto(posArgs[2])
	if err != nil {
		return err
	}

	weight, err := strconv.ParseUint(posArgs[4], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid weight: %v", err)
	}

	real := client.Real{
		Address: posArgs[3],
		Weight:  uint32(weight),
		Flags:   uint8(*flags),
	}

	if err := c.AddReal(posArgs[0], uint16(vipPort), vipProto, real); err != nil {
		return err
	}

	fmt.Printf("Real %s (weight=%d) added to VIP %s:%d/%s\n",
		real.Address, real.Weight, posArgs[0], vipPort, protoToString(vipProto))
	return nil
}

func realRemove(c *client.Client, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: real remove <vip-addr> <vip-port> <vip-proto> <real-addr>")
	}

	vipPort, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid VIP port: %v", err)
	}

	vipProto, err := parseProto(args[2])
	if err != nil {
		return err
	}

	if err := c.DeleteReal(args[0], uint16(vipPort), vipProto, args[3]); err != nil {
		return err
	}

	fmt.Printf("Real %s removed from VIP %s:%d/%s\n",
		args[3], args[0], vipPort, protoToString(vipProto))
	return nil
}

func realUpdate(c *client.Client, args []string) error {
	fs := flag.NewFlagSet("real update", flag.ExitOnError)
	flags := fs.Uint("flags", 0, "Real server flags")
	fs.Parse(args)

	posArgs := fs.Args()
	if len(posArgs) < 5 {
		return fmt.Errorf("usage: real update <vip-addr> <vip-port> <vip-proto> <real-addr> <weight> [--flags <flags>]")
	}

	vipPort, err := strconv.ParseUint(posArgs[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid VIP port: %v", err)
	}

	vipProto, err := parseProto(posArgs[2])
	if err != nil {
		return err
	}

	weight, err := strconv.ParseUint(posArgs[4], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid weight: %v", err)
	}

	reals := []client.Real{{
		Address: posArgs[3],
		Weight:  uint32(weight),
		Flags:   uint8(*flags),
	}}

	// Action 0 = add (which also updates existing)
	if err := c.UpdateReals(posArgs[0], uint16(vipPort), vipProto, 0, reals); err != nil {
		return err
	}

	fmt.Printf("Real %s updated (weight=%d) for VIP %s:%d/%s\n",
		posArgs[3], weight, posArgs[0], vipPort, protoToString(vipProto))
	return nil
}

// handleStats handles statistics commands.
func handleStats(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: stats <vip|lru|xdp|decap|ch-drop|icmp-too-big|src-routing|all> [args]")
	}

	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	watch := fs.Bool("watch", false, "Watch mode: show rate of change per second")
	interval := fs.Int("interval", 1, "Watch interval in seconds")

	// Find where flags start
	flagIdx := len(args)
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flagIdx = i
			break
		}
	}

	cmd := args[0]
	posArgs := args[1:flagIdx]
	fs.Parse(args[flagIdx:])

	if *watch {
		return statsWatch(c, cmd, posArgs, *interval)
	}

	return statsOnce(c, cmd, posArgs)
}

func statsOnce(c *client.Client, cmd string, args []string) error {
	switch cmd {
	case "vip":
		return statsVIP(c, args)
	case "lru":
		return statsLRU(c)
	case "xdp":
		return statsXDP(c)
	case "decap":
		return statsDecap(c)
	case "ch-drop":
		return statsCHDrop(c)
	case "icmp-too-big":
		return statsICMPTooBig(c)
	case "src-routing":
		return statsSrcRouting(c)
	case "all":
		return statsAll(c)
	default:
		return fmt.Errorf("unknown stats command: %s", cmd)
	}
}

func statsVIP(c *client.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: stats vip <address> <port> <proto>")
	}

	port, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	proto, err := parseProto(args[2])
	if err != nil {
		return err
	}

	stats, err := c.GetVIPStats(args[0], uint16(port), proto)
	if err != nil {
		return err
	}

	fmt.Printf("VIP Stats for %s:%d/%s:\n", args[0], port, protoToString(proto))
	fmt.Printf("  Packets: %d\n", stats.V1)
	fmt.Printf("  Bytes:   %d\n", stats.V2)
	return nil
}

func statsLRU(c *client.Client) error {
	stats, err := c.GetLRUStats()
	if err != nil {
		return err
	}
	miss, _ := c.GetLRUMissStats()
	fallback, _ := c.GetLRUFallbackStats()

	fmt.Println("LRU Stats:")
	fmt.Printf("  Total:    v1=%d v2=%d\n", stats.V1, stats.V2)
	if miss != nil {
		fmt.Printf("  Miss:     v1=%d v2=%d\n", miss.V1, miss.V2)
	}
	if fallback != nil {
		fmt.Printf("  Fallback: v1=%d v2=%d\n", fallback.V1, fallback.V2)
	}
	return nil
}

func statsXDP(c *client.Client) error {
	total, err := c.GetXDPTotalStats()
	if err != nil {
		return err
	}
	tx, _ := c.GetXDPTxStats()
	drop, _ := c.GetXDPDropStats()
	pass, _ := c.GetXDPPassStats()

	fmt.Println("XDP Stats:")
	fmt.Printf("  Total:   packets=%d bytes=%d\n", total.V1, total.V2)
	if tx != nil {
		fmt.Printf("  TX:      packets=%d bytes=%d\n", tx.V1, tx.V2)
	}
	if drop != nil {
		fmt.Printf("  Drop:    packets=%d bytes=%d\n", drop.V1, drop.V2)
	}
	if pass != nil {
		fmt.Printf("  Pass:    packets=%d bytes=%d\n", pass.V1, pass.V2)
	}
	return nil
}

func statsDecap(c *client.Client) error {
	stats, err := c.GetDecapStats()
	if err != nil {
		return err
	}
	inline, _ := c.GetInlineDecapStats()

	fmt.Println("Decap Stats:")
	fmt.Printf("  Decap:   packets=%d bytes=%d\n", stats.V1, stats.V2)
	if inline != nil {
		fmt.Printf("  Inline:  packets=%d bytes=%d\n", inline.V1, inline.V2)
	}
	return nil
}

func statsCHDrop(c *client.Client) error {
	stats, err := c.GetCHDropStats()
	if err != nil {
		return err
	}

	fmt.Println("CH Drop Stats:")
	fmt.Printf("  Packets: %d\n", stats.V1)
	fmt.Printf("  Bytes:   %d\n", stats.V2)
	return nil
}

func statsICMPTooBig(c *client.Client) error {
	stats, err := c.GetICMPTooBigStats()
	if err != nil {
		return err
	}

	fmt.Println("ICMP Too Big Stats:")
	fmt.Printf("  Packets: %d\n", stats.V1)
	fmt.Printf("  Bytes:   %d\n", stats.V2)
	return nil
}

func statsSrcRouting(c *client.Client) error {
	stats, err := c.GetSrcRoutingStats()
	if err != nil {
		return err
	}

	fmt.Println("Source Routing Stats:")
	fmt.Printf("  Packets: %d\n", stats.V1)
	fmt.Printf("  Bytes:   %d\n", stats.V2)
	return nil
}

func statsAll(c *client.Client) error {
	fmt.Println("=== All Statistics ===")
	fmt.Println()

	statsXDP(c)
	fmt.Println()
	statsLRU(c)
	fmt.Println()
	statsDecap(c)
	fmt.Println()
	statsCHDrop(c)
	fmt.Println()
	statsICMPTooBig(c)
	fmt.Println()
	statsSrcRouting(c)

	return nil
}

// statsEntry represents a statistics entry for watch mode.
type statsEntry struct {
	name string
	fn   func() (*client.LBStats, error)
}

// statsWatch displays stats in watch mode, showing rate of change.
func statsWatch(c *client.Client, cmd string, args []string, interval int) error {
	fmt.Printf("Watching %s stats (interval: %ds). Press Ctrl+C to stop.\n\n", cmd, interval)

	var entries []statsEntry

	switch cmd {
	case "vip":
		if len(args) < 3 {
			return fmt.Errorf("usage: stats vip <address> <port> <proto> --watch")
		}
		port, err := strconv.ParseUint(args[1], 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port: %v", err)
		}
		proto, err := parseProto(args[2])
		if err != nil {
			return err
		}
		entries = []statsEntry{
			{"VIP", func() (*client.LBStats, error) { return c.GetVIPStats(args[0], uint16(port), proto) }},
		}
	case "lru":
		entries = []statsEntry{
			{"LRU Total", c.GetLRUStats},
			{"LRU Miss", c.GetLRUMissStats},
			{"LRU Fallback", c.GetLRUFallbackStats},
		}
	case "xdp":
		entries = []statsEntry{
			{"XDP Total", c.GetXDPTotalStats},
			{"XDP TX", c.GetXDPTxStats},
			{"XDP Drop", c.GetXDPDropStats},
			{"XDP Pass", c.GetXDPPassStats},
		}
	case "decap":
		entries = []statsEntry{
			{"Decap", c.GetDecapStats},
			{"Inline Decap", c.GetInlineDecapStats},
		}
	case "ch-drop":
		entries = []statsEntry{
			{"CH Drop", c.GetCHDropStats},
		}
	case "icmp-too-big":
		entries = []statsEntry{
			{"ICMP Too Big", c.GetICMPTooBigStats},
		}
	case "src-routing":
		entries = []statsEntry{
			{"Src Routing", c.GetSrcRoutingStats},
		}
	case "all":
		entries = []statsEntry{
			{"XDP Total", c.GetXDPTotalStats},
			{"XDP TX", c.GetXDPTxStats},
			{"XDP Drop", c.GetXDPDropStats},
			{"XDP Pass", c.GetXDPPassStats},
			{"LRU Total", c.GetLRUStats},
			{"LRU Miss", c.GetLRUMissStats},
			{"Decap", c.GetDecapStats},
			{"CH Drop", c.GetCHDropStats},
		}
	default:
		return fmt.Errorf("watch not supported for: %s", cmd)
	}

	// Store previous values
	prev := make(map[string]*client.LBStats)

	// Print header
	fmt.Printf("%-20s %15s %15s %15s %15s\n", "COUNTER", "PACKETS", "PPS", "BYTES", "BPS")
	fmt.Println(strings.Repeat("-", 85))

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// First iteration
	printWatchStats(entries, prev, interval)

	for range ticker.C {
		// Move cursor up to overwrite
		fmt.Printf("\033[%dA", len(entries))
		printWatchStats(entries, prev, interval)
	}

	return nil
}

func printWatchStats(entries []statsEntry, prev map[string]*client.LBStats, interval int) {
	for _, e := range entries {
		stats, err := e.fn()
		if err != nil {
			fmt.Printf("%-20s %15s %15s %15s %15s\n", e.name, "error", "-", "error", "-")
			continue
		}

		var pps, bps int64
		if p, ok := prev[e.name]; ok {
			pps = int64(stats.V1-p.V1) / int64(interval)
			bps = int64(stats.V2-p.V2) / int64(interval)
		}
		prev[e.name] = stats

		fmt.Printf("%-20s %15d %15d %15d %15d\n", e.name, stats.V1, pps, stats.V2, bps)
	}
}

// handleMAC handles MAC address commands.
func handleMAC(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: mac <show|set> [mac-address]")
	}

	switch args[0] {
	case "show", "get":
		mac, err := c.GetMAC()
		if err != nil {
			return err
		}
		fmt.Printf("Default Router MAC: %s\n", mac)
		return nil
	case "set":
		if len(args) < 2 {
			return fmt.Errorf("usage: mac set <mac-address>")
		}
		if err := c.SetMAC(args[1]); err != nil {
			return err
		}
		fmt.Printf("MAC address set to: %s\n", args[1])
		return nil
	default:
		return fmt.Errorf("unknown mac command: %s", args[0])
	}
}

// handleHC handles healthcheck destination commands.
func handleHC(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: hc <list|add|remove> [args]")
	}

	switch args[0] {
	case "list", "ls":
		return hcList(c)
	case "add":
		return hcAdd(c, args[1:])
	case "remove", "rm", "delete", "del":
		return hcRemove(c, args[1:])
	default:
		return fmt.Errorf("unknown hc command: %s", args[0])
	}
}

func hcList(c *client.Client) error {
	dsts, err := c.ListHealthcheckerDsts()
	if err != nil {
		return err
	}

	if len(dsts) == 0 {
		fmt.Println("No healthcheck destinations configured")
		return nil
	}

	fmt.Printf("%-15s %-40s\n", "SOMARK", "DESTINATION")
	fmt.Println(strings.Repeat("-", 60))
	for _, dst := range dsts {
		fmt.Printf("%-15d %-40s\n", dst.Somark, dst.Dst)
	}
	return nil
}

func hcAdd(c *client.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: hc add <somark> <destination>")
	}

	somark, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid somark: %v", err)
	}

	if err := c.AddHealthcheckerDst(uint32(somark), args[1]); err != nil {
		return err
	}

	fmt.Printf("Healthcheck destination added: somark=%d dst=%s\n", somark, args[1])
	return nil
}

func hcRemove(c *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: hc remove <somark>")
	}

	somark, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid somark: %v", err)
	}

	if err := c.DeleteHealthcheckerDst(uint32(somark)); err != nil {
		return err
	}

	fmt.Printf("Healthcheck destination removed: somark=%d\n", somark)
	return nil
}

// handleHealth checks server health.
func handleHealth(c *client.Client) error {
	if err := c.Health(); err != nil {
		return fmt.Errorf("server unhealthy: %v", err)
	}
	fmt.Println("Server is healthy")
	return nil
}
