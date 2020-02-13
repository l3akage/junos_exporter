package policer

import (
	"strconv"

	"github.com/czerwonk/junos_exporter/collector"
	"github.com/czerwonk/junos_exporter/firewall"
	"github.com/czerwonk/junos_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "junos_policer_"

var (
	counterPackets *prometheus.Desc
	counterBytes   *prometheus.Desc
	policerPackets *prometheus.Desc
	policerBytes   *prometheus.Desc
)

func init() {
	l := []string{"target", "filter", "filter_group", "counter"}

	counterPackets = prometheus.NewDesc(prefix+"counter_packets", "Number of packets matching counter in firewall filter", l, nil)
	counterBytes = prometheus.NewDesc(prefix+"counter_bytes", "Number of bytes matching counter in firewall filter", l, nil)
	policerPackets = prometheus.NewDesc(prefix+"policer_packets", "Number of packets matching policer in firewall filter", l, nil)
	policerBytes = prometheus.NewDesc(prefix+"policer_bytes", "Number of bytes matching policer in firewall filter", l, nil)
}

type policerCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &policerCollector{}
}

// Name returns the name of the collector
func (*policerCollector) Name() string {
	return "Policer"
}

// Describe describes the metrics
func (*policerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- counterPackets
	ch <- counterBytes
	ch <- policerPackets
	ch <- policerBytes
}

// Collect collects metrics from JunOS
func (c *policerCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	var x = firewall.FirewallRpc{}
	err := client.RunCommandAndParse("show policer", &x)
	if err != nil {
		return err
	}
	cnt := 0
	for _, t := range x.Information.Filters {
		c.collectForFilter(cnt, t, ch, labelValues)
		cnt++
	}

	return nil
}

func (c *policerCollector) collectForFilter(groupID int, filter firewall.Filter, ch chan<- prometheus.Metric, labelValues []string) {
	l := append(labelValues, filter.Name, strconv.Itoa(groupID))

	for _, counter := range filter.Counters {
		lp := append(l, counter.Name)
		ch <- prometheus.MustNewConstMetric(counterPackets, prometheus.GaugeValue, float64(counter.Packets), lp...)
		ch <- prometheus.MustNewConstMetric(counterBytes, prometheus.GaugeValue, float64(counter.Bytes), lp...)
	}

	for _, policer := range filter.Policers {
		lp := append(l, policer.Name)
		ch <- prometheus.MustNewConstMetric(policerPackets, prometheus.GaugeValue, float64(policer.Packets), lp...)
		ch <- prometheus.MustNewConstMetric(policerBytes, prometheus.GaugeValue, float64(policer.Bytes), lp...)
	}
}
