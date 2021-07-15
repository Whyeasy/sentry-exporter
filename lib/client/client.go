package client

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sentry-exporter/internal"
	"github.com/whyeasy/sentry-exporter/lib/sentry"
)

//Stats struct is the list of expected results to export
type Stats struct {
	Issues *[]IssuesStats
}

//ExporterClient to build the client for the exporter.
type ExporterClient struct {
	client   *sentry.Client
	interval time.Duration
}

//New returns a new Client for connecting to Sentry and specifies the interval to renew cache.
func New(c internal.Config) *ExporterClient {

	convertedTime, _ := strconv.ParseInt(c.Interval, 10, 64)

	exporter := &ExporterClient{
		client:   sentry.NewClient(c.SentryAPIKey, c.SentryURI, c.SentryOrg),
		interval: time.Duration(convertedTime),
	}

	exporter.startFetchData()

	return exporter
}

var CachedStats *Stats = &Stats{
	Issues: &[]IssuesStats{},
}

//GetStats retrieves data from API to create metrics from.
func (c *ExporterClient) GetStats() (*Stats, error) {

	return CachedStats, nil
}

func (c *ExporterClient) getData() error {

	log.Info("Getting Data")
	projects, err := getProjects(c)
	if err != nil {
		return err
	}

	log.Info("Getting Issues")
	issueStats, err := getIssues(c, projects)
	if err != nil {
		return err
	}

	CachedStats = &Stats{
		Issues: issueStats,
	}

	log.Info("New data retrieved")

	return nil
}

func (c *ExporterClient) startFetchData() {

	// Do initial call to have data from the start.
	go func() {
		err := c.getData()
		if err != nil {
			log.Error("Scraping failed.")
		}
	}()

	ticker := time.NewTicker(c.interval * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.getData()
				if err != nil {
					log.Error("Scraping failed.")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
