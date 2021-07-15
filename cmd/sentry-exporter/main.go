package main

import (
	"flag"
	"fmt"
	"strings"

	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sentry-exporter/internal"
	"github.com/whyeasy/sentry-exporter/lib/client"
	"github.com/whyeasy/sentry-exporter/lib/collector"
)

var (
	config internal.Config
)

func init() {
	flag.StringVar(&config.LogFormat, "logFormat", os.Getenv("LOG_FORMAT"), "Default is logfmt, can be set to JSON.")
	flag.StringVar(&config.LogLevel, "logLevel", os.Getenv("LOG_LEVEL"), "Set different log level, default is Info.")
	flag.StringVar(&config.ListenAddress, "listenAddress", os.Getenv("LISTEN_ADDRESS"), "Port address of exporter to run on")
	flag.StringVar(&config.ListenPath, "listenPath", os.Getenv("LISTEN_PATH"), "Path where metrics will be exposed")
	flag.StringVar(&config.SentryURI, "sentryURI", os.Getenv("SENTRY_URI"), "URI to Sentry instance to monitor")
	flag.StringVar(&config.SentryAPIKey, "sentryAPIKey", os.Getenv("SENTRY_API_KEY"), "API Key to access Sentry")
	flag.StringVar(&config.SentryOrg, "sentryOrg", os.Getenv("SENTRY_ORG"), "Organization to scan in Sentry")
	flag.StringVar(&config.Interval, "interval", os.Getenv("INTERVAL"), "Provide a interval for scraping fresh data in seconds.")
}

func main() {
	if err := parseConfig(); err != nil {
		log.Error(err)
		flag.Usage()
		os.Exit(2)
	}
	initLogger()

	log.Info("Starting Sentry Exporter")

	client := client.New(config)
	coll := collector.New(client)
	prometheus.MustRegister(coll)

	log.Info("Start serving metrics")

	http.Handle(config.ListenPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>Sentry Exporter</title></head>
			<body>
			<h1>Sentry Exporter</h1>
			<p><a href="` + config.ListenPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Error(err)
		}
	})
	log.Fatal(http.ListenAndServe(":"+config.ListenAddress, nil))
}

func parseConfig() error {
	flag.Parse()
	required := []string{"sentryAPIKey", "sentryOrg"}
	var err error
	flag.VisitAll(func(f *flag.Flag) {
		for _, r := range required {
			if r == f.Name && (f.Value.String() == "" || f.Value.String() == "0") {
				err = fmt.Errorf("%v is empty", f.Usage)
			}
		}
		if f.Name == "sentryURI" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("https://sentry.io/api/0/")
			if err != nil {
				log.Error(err)
			}
		}
		if f.Name == "listenAddress" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("8080")
			if err != nil {
				log.Error(err)
			}
		}
		if f.Name == "listenPath" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("/metrics")
			if err != nil {
				log.Error(err)
			}
		}
		if f.Name == "interval" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("60")
			if err != nil {
				log.Error(err)
			}
		}
	})
	return err
}

func initLogger() {
	if strings.EqualFold(config.LogFormat, "json") {
		log.SetFormatter(&log.JSONFormatter{})
	}
	ll, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}
