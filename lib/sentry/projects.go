package sentry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//ProjectListResponse is the struct to unmarshal Sentry Projects
type ProjectListResponse []struct {
	ID           string `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Organization struct {
		ID     string `json:"id"`
		Slug   string `json:"slug"`
		Status struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		Name string `json:"name"`
	} `json:"organization"`
}

//ListProjects requests all projects from sentry.
func (c *Client) ListProjects() (*ProjectListResponse, error) {

	url := fmt.Sprintf("%sprojects/", c.sentryURI)
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

	var projects ProjectListResponse
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	log.Debug(projects)
	log.Debug(len(projects))

	return &projects, nil

}
