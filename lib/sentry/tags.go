package sentry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//TagsStructResponse is the struct to unmarshal Sentry Tags of projects.
type TagsStructResponse []struct {
	Key      string    `json:"key"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	LastSeen time.Time `json:"lastSeen"`
}

//Listtags requests all tags from a project in Sentry
func (c *Client) ListTags(project string, tag string) (*TagsStructResponse, error) {

	url := fmt.Sprintf("%sprojects/%s/%s/tags/%s/values/", c.sentryURI, c.sentryOrg, project, tag)
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
	if resp.StatusCode == 404 {
		log.Error("Tag not found: ", tag, " in Project:", project)
		return &TagsStructResponse{}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags TagsStructResponse
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	log.Debug(tags)
	log.Debug(len(tags))

	return &tags, nil

}
