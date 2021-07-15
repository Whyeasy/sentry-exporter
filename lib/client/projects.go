package client

import log "github.com/sirupsen/logrus"

//ProjectStats is the struct that holds the data we want from Jira
type ProjectStats struct {
	ID   string
	Slug string
}

func getProjects(c *ExporterClient) (*[]ProjectStats, error) {

	var result []ProjectStats

	projects, err := c.client.ListProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range *projects {
		if project.Status == "active" {
			result = append(result, ProjectStats{
				ID:   project.ID,
				Slug: project.Slug,
			})
		}
	}

	log.Info("Amount of projects found: ", len(*projects))

	return &result, nil
}
