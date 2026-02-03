package metrics

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
)

const namespace = "katran"

// KatranCollector implements prometheus.Collector for Katran metrics.
// It fetches statistics on-demand during each Prometheus scrape.
type KatranCollector struct {
	manager *lb.Manager

	// LB Status
	lbInitialized *prometheus.Desc
	lbReady       *prometheus.Desc

	// VIP metrics
	vipPacketsTotal     *prometheus.Desc
	vipBytesTotal       *prometheus.Desc
	vipDecapPacketsTotal *prometheus.Desc

	// XDP metrics
	xdpTotalPackets *prometheus.Desc
	xdpTotalBytes   *prometheus.Desc
	xdpTXPackets    *prometheus.Desc
	xdpTXBytes      *prometheus.Desc
	xdpDropPackets  *prometheus.Desc
	xdpDropBytes    *prometheus.Desc
	xdpPassPackets  *prometheus.Desc
	xdpPassBytes    *prometheus.Desc

	// LRU metrics
	lruTotalPackets    *prometheus.Desc
	lruHits            *prometheus.Desc
	lruMissTCPSyn      *prometheus.Desc
	lruMissNonSyn      *prometheus.Desc
	lruFallbackV1      *prometheus.Desc
	lruFallbackV2      *prometheus.Desc
	lruGlobalMapFailed *prometheus.Desc
	lruGlobalRouted    *prometheus.Desc

	// QUIC metrics
	quicCHRouted              *prometheus.Desc
	quicCIDInitial            *prometheus.Desc
	quicCIDInvalidServerID    *prometheus.Desc
	quicCIDInvalidServerIDSample *prometheus.Desc
	quicCIDRouted             *prometheus.Desc
	quicCIDUnknownRealDropped *prometheus.Desc
	quicCIDV0                 *prometheus.Desc
	quicCIDV1                 *prometheus.Desc
	quicCIDV2                 *prometheus.Desc
	quicCIDV3                 *prometheus.Desc
	quicDstMatchInLRU         *prometheus.Desc
	quicDstMismatchInLRU      *prometheus.Desc
	quicDstNotFoundInLRU      *prometheus.Desc

	// TPR metrics
	tprCHRouted         *prometheus.Desc
	tprSIDRouted        *prometheus.Desc
	tprDstMismatchInLRU *prometheus.Desc
	tprTCPSyn           *prometheus.Desc

	// Healthcheck metrics
	hcPacketsProcessed *prometheus.Desc
	hcPacketsDropped   *prometheus.Desc
	hcPacketsSkipped   *prometheus.Desc
	hcPacketsTooBig    *prometheus.Desc

	// ICMP metrics
	icmpTooBigV4 *prometheus.Desc
	icmpTooBigV6 *prometheus.Desc

	// CH Drop metrics
	chDropOutOfBounds *prometheus.Desc
	chDropUnmapped    *prometheus.Desc

	// Routing metrics
	srcRoutingLocal        *prometheus.Desc
	srcRoutingRemote       *prometheus.Desc
	inlineDecapV1          *prometheus.Desc
	inlineDecapV2          *prometheus.Desc
	decapV4                *prometheus.Desc
	decapV6                *prometheus.Desc

	// Userspace metrics
	bpfFailedCalls       *prometheus.Desc
	addrValidationFailed *prometheus.Desc

	// Per-core metrics
	perCorePackets *prometheus.Desc

	// Monitor metrics
	monitorLimit      *prometheus.Desc
	monitorAmount     *prometheus.Desc
	monitorBufferFull *prometheus.Desc

	// Flood metric
	underFlood *prometheus.Desc
}

// NewKatranCollector creates a new KatranCollector.
//
// Returns a new KatranCollector instance that can be registered with a Prometheus registry.
func NewKatranCollector() *KatranCollector {
	return &KatranCollector{
		manager: lb.GetManager(),

		// LB Status
		lbInitialized: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lb", "initialized"),
			"Whether the load balancer is initialized (1=yes, 0=no)",
			nil, nil,
		),
		lbReady: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lb", "ready"),
			"Whether the load balancer is ready (BPF programs loaded and attached)",
			nil, nil,
		),

		// VIP metrics
		vipPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "vip", "packets_total"),
			"Total packets processed for VIP",
			[]string{"address", "port", "proto"}, nil,
		),
		vipBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "vip", "bytes_total"),
			"Total bytes processed for VIP",
			[]string{"address", "port", "proto"}, nil,
		),
		vipDecapPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "vip", "decap_packets_total"),
			"Total decapsulated packets for VIP",
			[]string{"address", "port", "proto"}, nil,
		),

		// XDP metrics
		xdpTotalPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "total_packets"),
			"Total XDP packets processed",
			nil, nil,
		),
		xdpTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "total_bytes"),
			"Total XDP bytes processed",
			nil, nil,
		),
		xdpTXPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "tx_packets"),
			"XDP packets transmitted",
			nil, nil,
		),
		xdpTXBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "tx_bytes"),
			"XDP bytes transmitted",
			nil, nil,
		),
		xdpDropPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "drop_packets"),
			"XDP packets dropped",
			nil, nil,
		),
		xdpDropBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "drop_bytes"),
			"XDP bytes dropped",
			nil, nil,
		),
		xdpPassPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "pass_packets"),
			"XDP packets passed to kernel",
			nil, nil,
		),
		xdpPassBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "xdp", "pass_bytes"),
			"XDP bytes passed to kernel",
			nil, nil,
		),

		// LRU metrics
		lruTotalPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "total_packets"),
			"Total packets processed through LRU",
			nil, nil,
		),
		lruHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "hits"),
			"LRU cache hits",
			nil, nil,
		),
		lruMissTCPSyn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "miss_tcp_syn"),
			"LRU misses for TCP SYN packets",
			nil, nil,
		),
		lruMissNonSyn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "miss_non_syn"),
			"LRU misses for non-SYN packets",
			nil, nil,
		),
		lruFallbackV1: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "fallback_v1"),
			"LRU fallback V1 count",
			nil, nil,
		),
		lruFallbackV2: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "fallback_v2"),
			"LRU fallback V2 count",
			nil, nil,
		),
		lruGlobalMapFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "global_map_failed"),
			"Global LRU map lookup failures",
			nil, nil,
		),
		lruGlobalRouted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "lru", "global_routed"),
			"Packets routed via global LRU",
			nil, nil,
		),

		// QUIC metrics
		quicCHRouted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "ch_routed"),
			"QUIC packets routed via consistent hashing",
			nil, nil,
		),
		quicCIDInitial: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_initial"),
			"QUIC initial packets",
			nil, nil,
		),
		quicCIDInvalidServerID: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_invalid_server_id"),
			"QUIC packets with invalid server ID",
			nil, nil,
		),
		quicCIDInvalidServerIDSample: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_invalid_server_id_sample"),
			"QUIC packets with invalid server ID (sample)",
			nil, nil,
		),
		quicCIDRouted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_routed"),
			"QUIC packets routed via CID",
			nil, nil,
		),
		quicCIDUnknownRealDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_unknown_real_dropped"),
			"QUIC packets dropped due to unknown real",
			nil, nil,
		),
		quicCIDV0: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_v0"),
			"QUIC CID version 0 packets",
			nil, nil,
		),
		quicCIDV1: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_v1"),
			"QUIC CID version 1 packets",
			nil, nil,
		),
		quicCIDV2: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_v2"),
			"QUIC CID version 2 packets",
			nil, nil,
		),
		quicCIDV3: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "cid_v3"),
			"QUIC CID version 3 packets",
			nil, nil,
		),
		quicDstMatchInLRU: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "dst_match_in_lru"),
			"QUIC packets with destination match in LRU",
			nil, nil,
		),
		quicDstMismatchInLRU: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "dst_mismatch_in_lru"),
			"QUIC packets with destination mismatch in LRU",
			nil, nil,
		),
		quicDstNotFoundInLRU: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "quic", "dst_not_found_in_lru"),
			"QUIC packets with destination not found in LRU",
			nil, nil,
		),

		// TPR metrics
		tprCHRouted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tpr", "ch_routed"),
			"TPR packets routed via consistent hashing",
			nil, nil,
		),
		tprSIDRouted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tpr", "sid_routed"),
			"TPR packets routed via server ID",
			nil, nil,
		),
		tprDstMismatchInLRU: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tpr", "dst_mismatch_in_lru"),
			"TPR packets with destination mismatch in LRU",
			nil, nil,
		),
		tprTCPSyn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tpr", "tcp_syn"),
			"TPR TCP SYN packets processed",
			nil, nil,
		),

		// Healthcheck metrics
		hcPacketsProcessed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hc", "packets_processed"),
			"Healthcheck packets processed",
			nil, nil,
		),
		hcPacketsDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hc", "packets_dropped"),
			"Healthcheck packets dropped",
			nil, nil,
		),
		hcPacketsSkipped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hc", "packets_skipped"),
			"Healthcheck packets skipped",
			nil, nil,
		),
		hcPacketsTooBig: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hc", "packets_too_big"),
			"Healthcheck packets too big",
			nil, nil,
		),

		// ICMP metrics
		icmpTooBigV4: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "icmp", "too_big_v4"),
			"ICMPv4 too big messages",
			nil, nil,
		),
		icmpTooBigV6: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "icmp", "too_big_v6"),
			"ICMPv6 too big messages",
			nil, nil,
		),

		// CH Drop metrics
		chDropOutOfBounds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "ch_drop", "out_of_bounds"),
			"Consistent hash drops due to real ID out of bounds",
			nil, nil,
		),
		chDropUnmapped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "ch_drop", "unmapped"),
			"Consistent hash drops due to unmapped real",
			nil, nil,
		),

		// Routing metrics
		srcRoutingLocal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "src_routing", "local"),
			"Source routing to local backend",
			nil, nil,
		),
		srcRoutingRemote: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "src_routing", "remote"),
			"Source routing to remote backend (LPM matched)",
			nil, nil,
		),
		inlineDecapV1: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "inline_decap", "v1"),
			"Inline decap V1 count",
			nil, nil,
		),
		inlineDecapV2: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "inline_decap", "v2"),
			"Inline decap V2 count",
			nil, nil,
		),
		decapV4: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "decap", "v4"),
			"IPv4 packets decapsulated",
			nil, nil,
		),
		decapV6: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "decap", "v6"),
			"IPv6 packets decapsulated",
			nil, nil,
		),

		// Userspace metrics
		bpfFailedCalls: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "userspace", "bpf_failed_calls"),
			"Failed BPF syscalls",
			nil, nil,
		),
		addrValidationFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "userspace", "addr_validation_failed"),
			"Address validation failures",
			nil, nil,
		),

		// Per-core metrics
		perCorePackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "per_core", "packets"),
			"Packets processed per CPU core",
			[]string{"core"}, nil,
		),

		// Monitor metrics
		monitorLimit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "monitor", "limit"),
			"Monitor packet capture limit",
			nil, nil,
		),
		monitorAmount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "monitor", "amount"),
			"Monitor packets captured",
			nil, nil,
		),
		monitorBufferFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "monitor", "buffer_full"),
			"Monitor buffer full count",
			nil, nil,
		),

		// Flood metric
		underFlood: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "flood", "under_flood"),
			"Whether the system is under flood conditions (1=yes, 0=no)",
			nil, nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel.
func (c *KatranCollector) Describe(ch chan<- *prometheus.Desc) {
	// LB Status
	ch <- c.lbInitialized
	ch <- c.lbReady

	// VIP metrics
	ch <- c.vipPacketsTotal
	ch <- c.vipBytesTotal
	ch <- c.vipDecapPacketsTotal

	// XDP metrics
	ch <- c.xdpTotalPackets
	ch <- c.xdpTotalBytes
	ch <- c.xdpTXPackets
	ch <- c.xdpTXBytes
	ch <- c.xdpDropPackets
	ch <- c.xdpDropBytes
	ch <- c.xdpPassPackets
	ch <- c.xdpPassBytes

	// LRU metrics
	ch <- c.lruTotalPackets
	ch <- c.lruHits
	ch <- c.lruMissTCPSyn
	ch <- c.lruMissNonSyn
	ch <- c.lruFallbackV1
	ch <- c.lruFallbackV2
	ch <- c.lruGlobalMapFailed
	ch <- c.lruGlobalRouted

	// QUIC metrics
	ch <- c.quicCHRouted
	ch <- c.quicCIDInitial
	ch <- c.quicCIDInvalidServerID
	ch <- c.quicCIDInvalidServerIDSample
	ch <- c.quicCIDRouted
	ch <- c.quicCIDUnknownRealDropped
	ch <- c.quicCIDV0
	ch <- c.quicCIDV1
	ch <- c.quicCIDV2
	ch <- c.quicCIDV3
	ch <- c.quicDstMatchInLRU
	ch <- c.quicDstMismatchInLRU
	ch <- c.quicDstNotFoundInLRU

	// TPR metrics
	ch <- c.tprCHRouted
	ch <- c.tprSIDRouted
	ch <- c.tprDstMismatchInLRU
	ch <- c.tprTCPSyn

	// Healthcheck metrics
	ch <- c.hcPacketsProcessed
	ch <- c.hcPacketsDropped
	ch <- c.hcPacketsSkipped
	ch <- c.hcPacketsTooBig

	// ICMP metrics
	ch <- c.icmpTooBigV4
	ch <- c.icmpTooBigV6

	// CH Drop metrics
	ch <- c.chDropOutOfBounds
	ch <- c.chDropUnmapped

	// Routing metrics
	ch <- c.srcRoutingLocal
	ch <- c.srcRoutingRemote
	ch <- c.inlineDecapV1
	ch <- c.inlineDecapV2
	ch <- c.decapV4
	ch <- c.decapV6

	// Userspace metrics
	ch <- c.bpfFailedCalls
	ch <- c.addrValidationFailed

	// Per-core metrics
	ch <- c.perCorePackets

	// Monitor metrics
	ch <- c.monitorLimit
	ch <- c.monitorAmount
	ch <- c.monitorBufferFull

	// Flood metric
	ch <- c.underFlood
}

// Collect fetches the metrics from Katran and delivers them as Prometheus metrics.
func (c *KatranCollector) Collect(ch chan<- prometheus.Metric) {
	initialized, ready := c.manager.Status()

	// Always emit status metrics
	ch <- prometheus.MustNewConstMetric(c.lbInitialized, prometheus.GaugeValue, boolToFloat(initialized))
	ch <- prometheus.MustNewConstMetric(c.lbReady, prometheus.GaugeValue, boolToFloat(ready))

	// If LB is not initialized, we can't collect other metrics
	if !initialized {
		return
	}

	lbInstance, ok := c.manager.Get()
	if !ok {
		return
	}

	c.collectVIPMetrics(ch, lbInstance)
	c.collectXDPMetrics(ch, lbInstance)
	c.collectLRUMetrics(ch, lbInstance)
	c.collectQuicMetrics(ch, lbInstance)
	c.collectTPRMetrics(ch, lbInstance)
	c.collectHealthcheckMetrics(ch, lbInstance)
	c.collectICMPMetrics(ch, lbInstance)
	c.collectCHDropMetrics(ch, lbInstance)
	c.collectRoutingMetrics(ch, lbInstance)
	c.collectUserspaceMetrics(ch, lbInstance)
	c.collectPerCoreMetrics(ch, lbInstance)
	c.collectMonitorMetrics(ch, lbInstance)
	c.collectFloodMetrics(ch, lbInstance)
}

func (c *KatranCollector) collectVIPMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	vips, err := lb.GetAllVIPs()
	if err != nil {
		return
	}

	for _, vip := range vips {
		labels := []string{vip.Address, strconv.Itoa(int(vip.Port)), protoToString(vip.Proto)}

		// VIP stats
		if stats, err := lb.GetStatsForVIP(vip); err == nil {
			ch <- prometheus.MustNewConstMetric(c.vipPacketsTotal, prometheus.CounterValue, float64(stats.V1), labels...)
			ch <- prometheus.MustNewConstMetric(c.vipBytesTotal, prometheus.CounterValue, float64(stats.V2), labels...)
		}

		// Decap stats
		if stats, err := lb.GetDecapStatsForVIP(vip); err == nil {
			ch <- prometheus.MustNewConstMetric(c.vipDecapPacketsTotal, prometheus.CounterValue, float64(stats.V1), labels...)
		}
	}
}

func (c *KatranCollector) collectXDPMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	if stats, err := lb.GetXDPTotalStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.xdpTotalPackets, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.xdpTotalBytes, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetXDPTXStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.xdpTXPackets, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.xdpTXBytes, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetXDPDropStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.xdpDropPackets, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.xdpDropBytes, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetXDPPassStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.xdpPassPackets, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.xdpPassBytes, prometheus.CounterValue, float64(stats.V2))
	}
}

func (c *KatranCollector) collectLRUMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	if stats, err := lb.GetLRUStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.lruTotalPackets, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.lruHits, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetLRUMissStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.lruMissTCPSyn, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.lruMissNonSyn, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetLRUFallbackStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.lruFallbackV1, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.lruFallbackV2, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetGlobalLRUStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.lruGlobalMapFailed, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.lruGlobalRouted, prometheus.CounterValue, float64(stats.V2))
	}
}

func (c *KatranCollector) collectQuicMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetQuicPacketsStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.quicCHRouted, prometheus.CounterValue, float64(stats.CHRouted))
	ch <- prometheus.MustNewConstMetric(c.quicCIDInitial, prometheus.CounterValue, float64(stats.CIDInitial))
	ch <- prometheus.MustNewConstMetric(c.quicCIDInvalidServerID, prometheus.CounterValue, float64(stats.CIDInvalidServerID))
	ch <- prometheus.MustNewConstMetric(c.quicCIDInvalidServerIDSample, prometheus.CounterValue, float64(stats.CIDInvalidServerIDSample))
	ch <- prometheus.MustNewConstMetric(c.quicCIDRouted, prometheus.CounterValue, float64(stats.CIDRouted))
	ch <- prometheus.MustNewConstMetric(c.quicCIDUnknownRealDropped, prometheus.CounterValue, float64(stats.CIDUnknownRealDropped))
	ch <- prometheus.MustNewConstMetric(c.quicCIDV0, prometheus.CounterValue, float64(stats.CIDV0))
	ch <- prometheus.MustNewConstMetric(c.quicCIDV1, prometheus.CounterValue, float64(stats.CIDV1))
	ch <- prometheus.MustNewConstMetric(c.quicCIDV2, prometheus.CounterValue, float64(stats.CIDV2))
	ch <- prometheus.MustNewConstMetric(c.quicCIDV3, prometheus.CounterValue, float64(stats.CIDV3))
	ch <- prometheus.MustNewConstMetric(c.quicDstMatchInLRU, prometheus.CounterValue, float64(stats.DstMatchInLRU))
	ch <- prometheus.MustNewConstMetric(c.quicDstMismatchInLRU, prometheus.CounterValue, float64(stats.DstMismatchInLRU))
	ch <- prometheus.MustNewConstMetric(c.quicDstNotFoundInLRU, prometheus.CounterValue, float64(stats.DstNotFoundInLRU))
}

func (c *KatranCollector) collectTPRMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetTCPServerIDRoutingStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.tprCHRouted, prometheus.CounterValue, float64(stats.CHRouted))
	ch <- prometheus.MustNewConstMetric(c.tprSIDRouted, prometheus.CounterValue, float64(stats.SIDRouted))
	ch <- prometheus.MustNewConstMetric(c.tprDstMismatchInLRU, prometheus.CounterValue, float64(stats.DstMismatchInLRU))
	ch <- prometheus.MustNewConstMetric(c.tprTCPSyn, prometheus.CounterValue, float64(stats.TCPSyn))
}

func (c *KatranCollector) collectHealthcheckMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetHCProgStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.hcPacketsProcessed, prometheus.CounterValue, float64(stats.PacketsProcessed))
	ch <- prometheus.MustNewConstMetric(c.hcPacketsDropped, prometheus.CounterValue, float64(stats.PacketsDropped))
	ch <- prometheus.MustNewConstMetric(c.hcPacketsSkipped, prometheus.CounterValue, float64(stats.PacketsSkipped))
	ch <- prometheus.MustNewConstMetric(c.hcPacketsTooBig, prometheus.CounterValue, float64(stats.PacketsTooBig))
}

func (c *KatranCollector) collectICMPMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetICMPTooBigStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.icmpTooBigV4, prometheus.CounterValue, float64(stats.V1))
	ch <- prometheus.MustNewConstMetric(c.icmpTooBigV6, prometheus.CounterValue, float64(stats.V2))
}

func (c *KatranCollector) collectCHDropMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetCHDropStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.chDropOutOfBounds, prometheus.CounterValue, float64(stats.V1))
	ch <- prometheus.MustNewConstMetric(c.chDropUnmapped, prometheus.CounterValue, float64(stats.V2))
}

func (c *KatranCollector) collectRoutingMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	if stats, err := lb.GetSrcRoutingStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.srcRoutingLocal, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.srcRoutingRemote, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetInlineDecapStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.inlineDecapV1, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.inlineDecapV2, prometheus.CounterValue, float64(stats.V2))
	}

	if stats, err := lb.GetDecapStats(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.decapV4, prometheus.CounterValue, float64(stats.V1))
		ch <- prometheus.MustNewConstMetric(c.decapV6, prometheus.CounterValue, float64(stats.V2))
	}
}

func (c *KatranCollector) collectUserspaceMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetUserspaceStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.bpfFailedCalls, prometheus.CounterValue, float64(stats.BPFFailedCalls))
	ch <- prometheus.MustNewConstMetric(c.addrValidationFailed, prometheus.CounterValue, float64(stats.AddrValidationFailed))
}

func (c *KatranCollector) collectPerCoreMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	counts, err := lb.GetPerCorePacketsStats()
	if err != nil {
		return
	}

	for i, count := range counts {
		ch <- prometheus.MustNewConstMetric(c.perCorePackets, prometheus.CounterValue, float64(count), fmt.Sprintf("%d", i))
	}
}

func (c *KatranCollector) collectMonitorMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	stats, err := lb.GetMonitorStats()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.monitorLimit, prometheus.GaugeValue, float64(stats.Limit))
	ch <- prometheus.MustNewConstMetric(c.monitorAmount, prometheus.CounterValue, float64(stats.Amount))
	ch <- prometheus.MustNewConstMetric(c.monitorBufferFull, prometheus.CounterValue, float64(stats.BufferFull))
}

func (c *KatranCollector) collectFloodMetrics(ch chan<- prometheus.Metric, lb *katran.LoadBalancer) {
	underFlood, err := lb.IsUnderFlood()
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.underFlood, prometheus.GaugeValue, boolToFloat(underFlood))
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func protoToString(proto uint8) string {
	switch proto {
	case 6:
		return "tcp"
	case 17:
		return "udp"
	default:
		return fmt.Sprintf("%d", proto)
	}
}
