![build](https://github.com/Whyeasy/sentry-exporter/workflows/build/badge.svg)
![status-badge](https://goreportcard.com/badge/github.com/Whyeasy/sentry-exporter)
![Github go.mod Go version](https://img.shields.io/github/go-mod/go-version/Whyeasy/sentry-exporter)

# Sentry exporter

A Prometheus exporter for Sentry issues and events.

Currently this exporter retrieves the following metrics:

- Total amount of issues per project per environment
- Total amount of events per issue per project per environment

## Requirements

- Provide a Sentry API key; `--sentryAPIKey=<string>` or as env variable `SENTRY_API_KEY`.
- Provide your Sentry organization; `--sentryOrg=<string>` or as env variable `SENTRY_ORG`.

### Optional

Change listening port of the exporter; `--listenAddress <string>` or as env variable `LISTEN_ADDRESS`. Default = `8080`

Change listening path of the exporter; `--listenPath <string>` or as env variable `LISTEN_PATH`. Default = `/metrics`

Change the interval of retrieving data in the background; `--interval <string>` or as env variable `INTERVAL`. Default is `60`

Change the log format; `--logFormat=<string>` or as env variable `LOG_FORMAT`. Default = `logfmt`.

Change the log level; `--logLevel=<string>` or as env variable `LOG_LEVEL`. Default = `info`.

Change Sentry URI; `--sentryURI=<string>` or as env variable `SENTRY_URI`. Default = `https://sentry.io/api/0/`
