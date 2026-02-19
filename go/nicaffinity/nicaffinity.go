// Package nicaffinity provides NIC RX queue IRQ affinity management.
//
// It discovers RX queues for a network interface, maps their IRQs to CPUs
// with NUMA-aware placement, and optionally writes the affinity settings.
package nicaffinity

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// NUMATopology holds the mapping between CPUs and NUMA nodes.
type NUMATopology struct {
	// CPUToNode maps each CPU ID to its NUMA node ID.
	CPUToNode map[int]int
	// NodeToCPUs maps each NUMA node ID to its sorted list of CPU IDs.
	NodeToCPUs map[int][]int
}

// AffinityResult contains the CPU and NUMA node assignments from affinitization.
type AffinityResult struct {
	// CPUs is the sorted, deduplicated list of CPUs assigned to queues.
	CPUs []int
	// NUMANodes is the NUMA node ID for each corresponding CPU in CPUs.
	NUMANodes []int
}

// DiscoverRXQueues counts the number of RX queue directories under
// /sys/class/net/<iface>/queues/.
//
// iface: the network interface name to inspect.
func DiscoverRXQueues(iface string) (int, error) {
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

// GetNICNumaNode reads the NUMA node for a network interface from
// /sys/class/net/<iface>/device/numa_node. Returns -1 if the file
// contains -1 (common on single-node systems).
//
// iface: the network interface name.
func GetNICNumaNode(iface string) (int, error) {
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

// ReadNumaTopology reads /sys/devices/system/node/node*/cpulist to build
// a mapping between CPUs and NUMA nodes. Returns nil topology if the
// node directory cannot be read.
func ReadNumaTopology() (*NUMATopology, error) {
	nodeBase := "/sys/devices/system/node"
	entries, err := os.ReadDir(nodeBase)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", nodeBase, err)
	}

	topo := &NUMATopology{
		CPUToNode:  make(map[int]int),
		NodeToCPUs: make(map[int][]int),
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
		cpus := ParseCPUList(strings.TrimSpace(string(data)))
		sort.Ints(cpus)
		topo.NodeToCPUs[nodeID] = cpus
		for _, cpu := range cpus {
			topo.CPUToNode[cpu] = nodeID
		}
	}

	if len(topo.CPUToNode) == 0 {
		return nil, fmt.Errorf("no NUMA nodes found")
	}
	return topo, nil
}

// FindQueueIRQs parses /proc/interrupts to find IRQ numbers for each RX
// queue of the given interface. It recognizes three naming patterns:
// "<iface>-TxRx-N", "<iface>-rx-N", and "<iface>-N".
//
// iface: the network interface name.
// nQueues: the expected number of RX queues.
func FindQueueIRQs(iface string, nQueues int) (map[int]int, error) {
	f, err := os.Open("/proc/interrupts")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	patterns := []*regexp.Regexp{
		regexp.MustCompile(fmt.Sprintf(`%s-TxRx-(\d+)`, regexp.QuoteMeta(iface))),
		regexp.MustCompile(fmt.Sprintf(`%s-rx-(\d+)`, regexp.QuoteMeta(iface))),
		regexp.MustCompile(fmt.Sprintf(`%s-(\d+)`, regexp.QuoteMeta(iface))),
	}

	irqs := make(map[int]int)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		irqStr := strings.TrimSuffix(fields[0], ":")
		irqNum, err := strconv.Atoi(irqStr)
		if err != nil {
			continue
		}

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

// SelectCPUs builds an ordered list of CPUs for queue assignment. CPUs on
// the NIC's local NUMA node come first, followed by CPUs on remote nodes
// (sorted by node ID then CPU ID). If nicNode is -1, all CPUs are treated
// as local.
//
// n: number of queues (informational, not used for slicing).
// nicNode: the NIC's NUMA node ID, or -1 if unknown.
// topo: the NUMA topology, or nil to fall back to online CPUs.
func SelectCPUs(n int, nicNode int, topo *NUMATopology) []int {
	if topo == nil || len(topo.CPUToNode) == 0 {
		return FallbackCPUs()
	}

	if nicNode == -1 {
		var all []int
		for cpu := range topo.CPUToNode {
			all = append(all, cpu)
		}
		sort.Ints(all)
		return all
	}

	local := make([]int, len(topo.NodeToCPUs[nicNode]))
	copy(local, topo.NodeToCPUs[nicNode])
	sort.Ints(local)

	var remoteNodes []int
	for nodeID := range topo.NodeToCPUs {
		if nodeID != nicNode {
			remoteNodes = append(remoteNodes, nodeID)
		}
	}
	sort.Ints(remoteNodes)

	var remote []int
	for _, nodeID := range remoteNodes {
		cpus := make([]int, len(topo.NodeToCPUs[nodeID]))
		copy(cpus, topo.NodeToCPUs[nodeID])
		sort.Ints(cpus)
		remote = append(remote, cpus...)
	}

	return append(local, remote...)
}

// SetIRQAffinity writes a single CPU ID to /proc/irq/<irq>/smp_affinity_list
// to pin the given IRQ to that CPU.
//
// irq: the IRQ number to configure.
// cpu: the CPU ID to assign.
func SetIRQAffinity(irq, cpu int) error {
	path := fmt.Sprintf("/proc/irq/%d/smp_affinity_list", irq)
	return os.WriteFile(path, []byte(strconv.Itoa(cpu)), 0644)
}

// ParseCPUList parses a CPU list string such as "0-3,8-11" into a slice
// of individual CPU IDs. It handles both single values ("5") and ranges
// ("0-3").
//
// s: the CPU list string to parse.
func ParseCPUList(s string) []int {
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

// FallbackCPUs reads /sys/devices/system/cpu/online to get available CPUs
// when NUMA topology is unavailable. Returns a sorted list of CPU IDs.
func FallbackCPUs() []int {
	data, err := os.ReadFile("/sys/devices/system/cpu/online")
	if err != nil {
		return nil
	}
	cpus := ParseCPUList(strings.TrimSpace(string(data)))
	sort.Ints(cpus)
	return cpus
}

// Affinitize is the main orchestrator. It discovers RX queues for the given
// interface, finds their IRQs, selects CPUs with NUMA awareness, and
// optionally writes affinity settings.
//
// iface: the network interface name.
// dryRun: if true, skip writing affinity (preview only).
//
// Returns an AffinityResult with unique CPUs and their NUMA nodes.
func Affinitize(iface string, dryRun bool) (*AffinityResult, error) {
	nQueues, err := DiscoverRXQueues(iface)
	if err != nil {
		return nil, fmt.Errorf("discovering RX queues: %w", err)
	}
	if nQueues == 0 {
		return nil, fmt.Errorf("no RX queues found for interface %s", iface)
	}

	nicNode, err := GetNICNumaNode(iface)
	if err != nil {
		nicNode = -1
	}

	topo, err := ReadNumaTopology()
	if err != nil {
		topo = nil
	}

	irqs, err := FindQueueIRQs(iface, nQueues)
	if err != nil {
		return nil, fmt.Errorf("finding queue IRQs: %w", err)
	}

	orderedCPUs := SelectCPUs(nQueues, nicNode, topo)
	if len(orderedCPUs) == 0 {
		return nil, fmt.Errorf("no CPUs available for assignment")
	}

	if !dryRun && os.Geteuid() != 0 {
		return nil, fmt.Errorf("writing IRQ affinity requires root privileges; use dry_run to preview")
	}

	// Track unique CPUs assigned (deduplicated, sorted).
	cpuSet := make(map[int]struct{})
	for i := 0; i < nQueues; i++ {
		cpu := orderedCPUs[i%len(orderedCPUs)]
		cpuSet[cpu] = struct{}{}

		irq, ok := irqs[i]
		if !ok {
			continue
		}

		if !dryRun {
			if err := SetIRQAffinity(irq, cpu); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to set affinity for IRQ %d (queue %d): %v\n", irq, i, err)
			}
		}
	}

	// Build sorted, deduplicated CPU list and corresponding NUMA nodes.
	var cpus []int
	for cpu := range cpuSet {
		cpus = append(cpus, cpu)
	}
	sort.Ints(cpus)

	numaNodes := make([]int, len(cpus))
	for i, cpu := range cpus {
		if topo != nil {
			if n, exists := topo.CPUToNode[cpu]; exists {
				numaNodes[i] = n
			}
		}
	}

	return &AffinityResult{
		CPUs:      cpus,
		NUMANodes: numaNodes,
	}, nil
}
