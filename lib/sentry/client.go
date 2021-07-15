package sentry

import (
	"net/http"
	"strings"
)

//Client struct holds the data we need for the Sentry Client
type Client struct {
	sentryURI    string
	client       *http.Client
	sentryAPIKey string
	sentryOrg    string
}

//NewClient Creates a new client to communicate with the Sentry API.
func NewClient(api string, baseURI string, org string) *Client {

	if !strings.HasSuffix(baseURI, "/") {
		baseURI += "/"
	}

	return &Client{
		client:       &http.Client{},
		sentryAPIKey: api,
		sentryURI:    baseURI,
		sentryOrg:    org,
	}
}
