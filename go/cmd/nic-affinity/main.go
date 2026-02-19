// nic-affinity configures NIC RX queue IRQ affinity across CPUs.
//
// It maps each RX queue's IRQ to a single CPU, preferring CPUs on the NIC's
// local NUMA node, and prints the CPU/NUMA mapping on exit.
//
// Usage:
//
//	nic-affinity -iface <name> [-dry-run]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// numaTopology holds the mapping between CPUs and NUMA nodes.
type numaTopology struct {
	// cpuToNode maps each CPU ID to its NUMA node ID.
	cpuToNode map[int]int
	// nodeToCPUs maps each NUMA node ID to its sorted list of CPU IDs.
	nodeToCPUs map[int][]int
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main logic. It discovers RX queues, finds their IRQs,
// selects CPUs with NUMA awareness, and writes affinity settings.
// Returns an error if any critical step fails.
func run() error {
	iface := flag.String("iface", "", "network interface name (required)")
	dryRun := flag.Bool("dry-run", false, "print assignments without writing affinity")
	flag.Parse()

	if *iface == "" {
		return fmt.Errorf("missing required flag: -iface")
	}

	nQueues, err := discoverRXQueues(*iface)
	if err != nil {
		return fmt.Errorf("discovering RX queues: %w", err)
	}
	if nQueues == 0 {
		return fmt.Errorf("no RX queues found for interface %s", *iface)
	}

	nicNode, err := getNICNumaNode(*iface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not read NIC NUMA node: %v; treating all CPUs as local\n", err)
		nicNode = -1
	}

	topo, err := readNumaTopology()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not read NUMA topology: %v; falling back to flat CPU list\n", err)
		topo = nil
	}

	irqs, err := findQueueIRQs(*iface, nQueues)
	if err != nil {
		return fmt.Errorf("finding queue IRQs: %w", err)
	}

	orderedCPUs := selectCPUs(nQueues, nicNode, topo)
	if len(orderedCPUs) == 0 {
		return fmt.Errorf("no CPUs available for assignment")
	}

	// Check root before attempting writes.
	if !*dryRun && os.Geteuid() != 0 {
		return fmt.Errorf("writing IRQ affinity requires root privileges; use -dry-run to preview")
	}

	cpuAssignments := make([]int, nQueues)
	for i := 0; i < nQueues; i++ {
		cpu := orderedCPUs[i%len(orderedCPUs)]
		cpuAssignments[i] = cpu

		irq, ok := irqs[i]
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: no IRQ found for queue %d, skipping\n", i)
			continue
		}

		if *dryRun {
			node := -1
			if topo != nil {
				if n, exists := topo.cpuToNode[cpu]; exists {
					node = n
				}
			}
			fmt.Fprintf(os.Stderr, "queue %d: IRQ %d -> CPU %d (NUMA %d)\n", i, irq, cpu, node)
		} else {
			if err := setIRQAffinity(irq, cpu); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to set affinity for IRQ %d (queue %d): %v\n", irq, i, err)
			}
		}
	}

	// Print output: line 1 = CPU IDs, line 2 = NUMA node IDs.
	cpuStrs := make([]string, nQueues)
	nodeStrs := make([]string, nQueues)
	for i, cpu := range cpuAssignments {
		cpuStrs[i] = strconv.Itoa(cpu)
		node := 0
		if topo != nil {
			if n, exists := topo.cpuToNode[cpu]; exists {
				node = n
			}
		}
		nodeStrs[i] = strconv.Itoa(node)
	}
	fmt.Println(strings.Join(cpuStrs, ","))
	fmt.Println(strings.Join(nodeStrs, ","))

	return nil
}

// discoverRXQueues counts the number of RX queue directories under
// /sys/class/net/<iface>/queues/. Returns the count and any error
// encountered while reading the directory.
//
// iface: the network interface name to inspect.
func discoverRXQueues(iface string) (int, error) {
	queueDir := filepath.Join("/sys/class/net", iface, "queues")
	entries, err := os.ReadDir(queueDir)
	if err != nil {
		return 0, fmt.Errorf("reading %s: %w", queueDir, err)
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "rx-") {
			count++
		}
	}
	return count, nil
}

// getNICNumaNode reads the NUMA node for a network interface from
// /sys/class/net/<iface>/device/numa_node. Returns -1 if the file
// contains -1 (common on single-node systems).
//
// iface: the network interface name.
func getNICNumaNode(iface string) (int, error) {
	path := filepath.Join("/sys/class/net", iface, "device", "numa_node")
	data, err := os.ReadFile(path)
	if err != nil {
		return -1, err
	}
	node, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return -1, fmt.Errorf("parsing numa_node: %w", err)
	}
	return node, nil
}

// readNumaTopology reads /sys/devices/system/node/node*/cpulist to build
// a mapping between CPUs and NUMA nodes. Returns nil topology if the
// node directory cannot be read.
func readNumaTopology() (*numaTopology, error) {
	nodeBase := "/sys/devices/system/node"
	entries, err := os.ReadDir(nodeBase)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", nodeBase, err)
	}

	topo := &numaTopology{
		cpuToNode:  make(map[int]int),
		nodeToCPUs: make(map[int][]int),
	}

	nodeRe := regexp.MustCompile(`^node(\d+)$`)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := nodeRe.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		nodeID, _ := strconv.Atoi(m[1])

		cpulistPath := filepath.Join(nodeBase, e.Name(), "cpulist")
		data, err := os.ReadFile(cpulistPath)
		if err != nil {
			continue
		}
		cpus := parseCPUList(strings.TrimSpace(string(data)))
		sort.Ints(cpus)
		topo.nodeToCPUs[nodeID] = cpus
		for _, cpu := range cpus {
			topo.cpuToNode[cpu] = nodeID
		}
	}

	if len(topo.cpuToNode) == 0 {
		return nil, fmt.Errorf("no NUMA nodes found")
	}
	return topo, nil
}

// findQueueIRQs parses /proc/interrupts to find IRQ numbers for each RX
// queue of the given interface. It recognizes three naming patterns:
// "<iface>-TxRx-N", "<iface>-rx-N", and "<iface>-N".
//
// iface: the network interface name.
// nQueues: the expected number of RX queues.
//
// Returns a map from queue index to IRQ number.
func findQueueIRQs(iface string, nQueues int) (map[int]int, error) {
	f, err := os.Open("/proc/interrupts")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Patterns to match queue IRQ names.
	patterns := []*regexp.Regexp{
		regexp.MustCompile(fmt.Sprintf(`%s-TxRx-(\d+)`, regexp.QuoteMeta(iface))),
		regexp.MustCompile(fmt.Sprintf(`%s-rx-(\d+)`, regexp.QuoteMeta(iface))),
		regexp.MustCompile(fmt.Sprintf(`%s-(\d+)`, regexp.QuoteMeta(iface))),
	}

	irqs := make(map[int]int) // queue index -> IRQ number
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		// First field is the IRQ number followed by a colon.
		irqStr := strings.TrimSuffix(fields[0], ":")
		irqNum, err := strconv.Atoi(irqStr)
		if err != nil {
			continue
		}

		// The last field is typically the device name.
		devName := fields[len(fields)-1]
		for _, pat := range patterns {
			m := pat.FindStringSubmatch(devName)
			if m == nil {
				continue
			}
			queueIdx, _ := strconv.Atoi(m[1])
			if queueIdx < nQueues {
				if _, exists := irqs[queueIdx]; !exists {
					irqs[queueIdx] = irqNum
				}
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading /proc/interrupts: %w", err)
	}
	return irqs, nil
}

// selectCPUs builds an ordered list of CPUs for queue assignment. CPUs on
// the NIC's local NUMA node come first, followed by CPUs on remote nodes
// (sorted by node ID then CPU ID). If nicNode is -1, all CPUs are treated
// as local.
//
// n: number of queues (informational, not used for slicing).
// nicNode: the NIC's NUMA node ID, or -1 if unknown.
// topo: the NUMA topology, or nil to fall back to online CPUs.
func selectCPUs(n int, nicNode int, topo *numaTopology) []int {
	if topo == nil || len(topo.cpuToNode) == 0 {
		return fallbackCPUs()
	}

	if nicNode == -1 {
		// All CPUs as local, sorted.
		var all []int
		for cpu := range topo.cpuToNode {
			all = append(all, cpu)
		}
		sort.Ints(all)
		return all
	}

	// Local CPUs first.
	local := make([]int, len(topo.nodeToCPUs[nicNode]))
	copy(local, topo.nodeToCPUs[nicNode])
	sort.Ints(local)

	// Remote CPUs: sorted by node then CPU.
	var remoteNodes []int
	for nodeID := range topo.nodeToCPUs {
		if nodeID != nicNode {
			remoteNodes = append(remoteNodes, nodeID)
		}
	}
	sort.Ints(remoteNodes)

	var remote []int
	for _, nodeID := range remoteNodes {
		cpus := make([]int, len(topo.nodeToCPUs[nodeID]))
		copy(cpus, topo.nodeToCPUs[nodeID])
		sort.Ints(cpus)
		remote = append(remote, cpus...)
	}

	return append(local, remote...)
}

// fallbackCPUs reads /sys/devices/system/cpu/online to get available CPUs
// when NUMA topology is unavailable. Returns a sorted list of CPU IDs.
func fallbackCPUs() []int {
	data, err := os.ReadFile("/sys/devices/system/cpu/online")
	if err != nil {
		return nil
	}
	cpus := parseCPUList(strings.TrimSpace(string(data)))
	sort.Ints(cpus)
	return cpus
}

// setIRQAffinity writes a single CPU ID to /proc/irq/<irq>/smp_affinity_list
// to pin the given IRQ to that CPU.
//
// irq: the IRQ number to configure.
// cpu: the CPU ID to assign.
func setIRQAffinity(irq, cpu int) error {
	path := fmt.Sprintf("/proc/irq/%d/smp_affinity_list", irq)
	return os.WriteFile(path, []byte(strconv.Itoa(cpu)), 0644)
}

// parseCPUList parses a CPU list string such as "0-3,8-11" into a slice
// of individual CPU IDs. It handles both single values ("5") and ranges
// ("0-3").
//
// s: the CPU list string to parse.
func parseCPUList(s string) []int {
	var result []int
	if s == "" {
		return result
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "-"); idx >= 0 {
			lo, err1 := strconv.Atoi(part[:idx])
			hi, err2 := strconv.Atoi(part[idx+1:])
			if err1 != nil || err2 != nil {
				continue
			}
			for i := lo; i <= hi; i++ {
				result = append(result, i)
			}
		} else {
			v, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			result = append(result, v)
		}
	}
	return result
}
