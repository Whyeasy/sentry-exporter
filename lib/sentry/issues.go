package sentry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//IssueListResponse is the struct to unmarshal Sentry issues
type IssueListResponse []struct {
	ID            string      `json:"id"`
	ShareID       interface{} `json:"shareId"`
	ShortID       string      `json:"shortId"`
	Title         string      `json:"title"`
	Culprit       string      `json:"culprit"`
	Permalink     string      `json:"permalink"`
	Logger        string      `json:"logger"`
	Level         string      `json:"level"`
	Status        string      `json:"status"`
	StatusDetails struct {
	} `json:"statusDetails"`
	Type                string        `json:"type"`
	NumComments         int           `json:"numComments"`
	IsBookmarked        bool          `json:"isBookmarked"`
	IsSubscribed        bool          `json:"isSubscribed"`
	SubscriptionDetails interface{}   `json:"subscriptionDetails"`
	HasSeen             bool          `json:"hasSeen"`
	Annotations         []interface{} `json:"annotations"`
	IsUnhandled         bool          `json:"isUnhandled"`
	Count               string        `json:"count"`
	UserCount           int           `json:"userCount"`
	FirstSeen           time.Time     `json:"firstSeen"`
	LastSeen            time.Time     `json:"lastSeen"`
	Stats               struct {
		Two4H [][]int `json:"24h"`
		One4D [][]int `json:"14d"`
	} `json:"stats"`
}

//Listissues requests all issues from sentry per project per environment in 24h and 14d measurements.
func (c *Client) ListIssues(project string, projectID string, env string, period string) (*IssueListResponse, error) {

	url := fmt.Sprintf("%sprojects/%s/%s/issues/?project=%s&sort=date&environment=%s&statsPeriod=%s&query=lastSeen:-%s", c.sentryURI, c.sentryOrg, project, projectID, env, period, period)
	log.Debug(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + c.sentryAPIKey
	req.Header.Add("Authorization", bearer)

	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var issues IssueListResponse
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, err
	}

	log.Debug(issues)
	log.Debug(len(issues))

	return &issues, nil

}
