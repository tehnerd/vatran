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
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tehnerd/vatran/go/nicaffinity"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main logic. It uses the nicaffinity library to discover
// RX queues, find their IRQs, select CPUs with NUMA awareness, and write
// affinity settings. In dry-run mode, it also prints per-queue assignments
// to stderr.
// Returns an error if any critical step fails.
func run() error {
	iface := flag.String("iface", "", "network interface name (required)")
	dryRun := flag.Bool("dry-run", false, "print assignments without writing affinity")
	flag.Parse()

	if *iface == "" {
		return fmt.Errorf("missing required flag: -iface")
	}

	// In dry-run mode, print per-queue details before the summary.
	if *dryRun {
		if err := printDryRunDetails(*iface); err != nil {
			// Non-fatal: fall through to Affinitize for the summary.
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		}
	}

	result, err := nicaffinity.Affinitize(*iface, *dryRun)
	if err != nil {
		return err
	}

	// Print output: line 1 = CPU IDs, line 2 = NUMA node IDs.
	cpuStrs := make([]string, len(result.CPUs))
	nodeStrs := make([]string, len(result.NUMANodes))
	for i, cpu := range result.CPUs {
		cpuStrs[i] = strconv.Itoa(cpu)
		nodeStrs[i] = strconv.Itoa(result.NUMANodes[i])
	}
	fmt.Println(strings.Join(cpuStrs, ","))
	fmt.Println(strings.Join(nodeStrs, ","))

	return nil
}

// printDryRunDetails prints per-queue IRQ-to-CPU assignments to stderr.
//
// iface: the network interface name.
func printDryRunDetails(iface string) error {
	nQueues, err := nicaffinity.DiscoverRXQueues(iface)
	if err != nil {
		return fmt.Errorf("discovering RX queues: %w", err)
	}

	nicNode, _ := nicaffinity.GetNICNumaNode(iface)
	topo, _ := nicaffinity.ReadNumaTopology()

	irqs, err := nicaffinity.FindQueueIRQs(iface, nQueues)
	if err != nil {
		return fmt.Errorf("finding queue IRQs: %w", err)
	}

	orderedCPUs := nicaffinity.SelectCPUs(nQueues, nicNode, topo)
	for i := 0; i < nQueues; i++ {
		cpu := orderedCPUs[i%len(orderedCPUs)]
		irq, ok := irqs[i]
		if !ok {
			fmt.Fprintf(os.Stderr, "queue %d: no IRQ found, skipping\n", i)
			continue
		}
		node := -1
		if topo != nil {
			if n, exists := topo.CPUToNode[cpu]; exists {
				node = n
			}
		}
		fmt.Fprintf(os.Stderr, "queue %d: IRQ %d -> CPU %d (NUMA %d)\n", i, irq, cpu, node)
	}
	return nil
}
