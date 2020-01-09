package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, webhook string, name string, commitUrl string) {
	buildUrl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s", b.Id)
	startTime, _ := time.Parse(time.RFC3339, b.StartTime)
	finishTime, _ := time.Parse(time.RFC3339, b.FinishTime)
	duration := finishTime.Sub(startTime)
	var i string
	switch b.Status {
	case "SUCCESS":
		i = ":white_check_mark:"
	case "FAILURE", "CANCELLED":
		i = ":x:"
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		i = ":interrobang:"
	default:
		i = ":question:"
	}
	j := fmt.Sprintf(
		`{
		    "blocks": [
				{
					"type": "section",
					"text": {
						"type": "plain_text",
						"text": "%s Cloud Build %s completed: %s %s in %s"
					}
				},
				{
					"type": "actions",
					"elements": [
						{
							"type": "button",
							"text": {
								"type": "plain_text",
								"text": "Open build details"
							},
							"url": "%s"
						},
						{
							"type": "button",
							"text": {
								"type": "plain_text",
								"text": "Open commit details"
							},
							"url": "%s"
						}
					]
				}
			]
		}`, name, b.Id, i, b.Status, duration, buildUrl, commitUrl)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
