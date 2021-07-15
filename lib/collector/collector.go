package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/whyeasy/sentry-exporter/lib/client"
)

//Collector struct for holding Prometheus Desc and Exporter Client
type Collector struct {
	up     *prometheus.Desc
	client *client.ExporterClient

	events *prometheus.Desc
	issues *prometheus.Desc
}

//New creates a new Collecotor with Prometheus descriptors
func New(c *client.ExporterClient) *Collector {
	log.Info("Creating collector")
	return &Collector{
		up:     prometheus.NewDesc("sentry_up", "Whether Sentry scrap was successful", nil, nil),
		client: c,

		events: prometheus.NewDesc("sentry_events", "Total events triggered divided per project and environment", []string{"project_name", "environment", "period"}, nil),
		issues: prometheus.NewDesc("sentry_issues", "Total issues happening in each project and environment ", []string{"project_name", "environment", "period"}, nil),
	}
}

//Describe the metrics that are collected.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.events
}

//Collect gathers the metrics that are exported.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Info("Running scrape")

	if stats, err := c.client.GetStats(); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectIssueStats(c, ch, stats)

		log.Info("Scrape Complete")
	}
}

func collectIssueStats(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, issue := range *stats.Issues {
		ch <- prometheus.MustNewConstMetric(c.events, prometheus.GaugeValue, float64(issue.EventsTotal), issue.Project, issue.Env, issue.Period)
		ch <- prometheus.MustNewConstMetric(c.issues, prometheus.GaugeValue, float64(issue.IssuesTotal), issue.Project, issue.Env, issue.Period)
	}
}
