package client

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//IssueStats is the struct that holds all the
type IssuesStats struct {
	Project     string
	Env         string
	IssueID     string
	IssuesTotal int
	EventsTotal int
	Period      string
}

func getIssues(c *ExporterClient, projects *[]ProjectStats) (*[]IssuesStats, error) {

	var result []IssuesStats
	periods := []string{"24h", "14d"}

	var wg sync.WaitGroup
	wg.Add(len(*projects))

	log.Debug("Projects: ", *projects)

	for _, project := range *projects {

		go func(projectName string, projectID string) {

			defer wg.Done()
			envs, err := c.client.ListTags(projectName, "environment")
			if err != nil {
				log.Fatal("Error: ", err, " ProjectID: ", projectID)
			}

			for _, env := range *envs {

				newerTags := time.Now().Add(-336 * time.Hour)

				if env.LastSeen.After(newerTags) {
					for _, period := range periods {

						issues, err := c.client.ListIssues(projectName, projectID, env.Value, period)
						if err != nil {
							log.Fatal(err)
						}

						var issuesEvents int
						for _, issue := range *issues {
							if period == "24h" {
								for _, stats := range issue.Stats.Two4H {
									issuesEvents += stats[1]
								}
							} else {
								for _, stats := range issue.Stats.One4D {
									issuesEvents += stats[1]
								}
							}
						}

						result = append(result, IssuesStats{
							Env:         env.Value,
							IssuesTotal: len(*issues),
							Project:     projectName,
							EventsTotal: issuesEvents,
							Period:      period,
						})

						log.Info("Amount of issues found: ", len(*issues), " total amount of events: ", issuesEvents, " for project: ", projectName, " env: ", env.Value, " period: ", period)
					}
				}
			}

		}(project.Slug, project.ID)
	}

	wg.Wait()

	return &result, nil
}
