package subscribers

import (
	"github.com/czerwonk/junos_exporter/collector"
	"github.com/czerwonk/junos_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "junos_subscribers_"

var (
	usernameInterfaceDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "username", "name"}
	usernameInterfaceDesc = prometheus.NewDesc(prefix+"interface", "Maps subscriber username to interface", l, nil)
}

type subscribersCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &subscribersCollector{}
}

// Name returns the name of the collector
func (*subscribersCollector) Name() string {
	return "Subscribers"
}

// Describe describes the metrics
func (*subscribersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- usernameInterfaceDesc
}

// Collect collects metrics from JunOS
func (c *subscribersCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	var x = SubscribersRPC{}
	err := client.RunCommandAndParse("show subscribers", &x)
	if err != nil {
		return err
	}

	for _, s := range x.SubscribersInformation.Subscribers {
		if len(s.Username) <= 0 {
			continue
		}
		l := append(labelValues, s.Username, s.Interface)
		ch <- prometheus.MustNewConstMetric(usernameInterfaceDesc, prometheus.GaugeValue, 1, l...)
	}

	return nil
}
