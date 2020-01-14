package slackbot

import (
	"context"
	"fmt"
	"log"

	// duration "github.com/golang/protobuf/ptypes/duration"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Trigger starts an independent watcher build.
func Trigger(ctx context.Context, build string, webhook string, name string, commitUrl string) {
	svc := gcbClient(ctx)
	project, err := getProject()
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	log.Printf("Getting build timeout %s", build)
	lc := svc.Projects.Builds.Get(project, build)
	thisBuild, err := lc.Do()

	log.Printf("Build timeout was: %s", thisBuild.Timeout)

	b := &cloudbuild.Build{
		Steps: []*cloudbuild.BuildStep{
			&cloudbuild.BuildStep{
				Name: "gcr.io/$PROJECT_ID/slackbot",
				Args: []string{
					fmt.Sprintf("--build=%s", build),
					fmt.Sprintf("--webhook=%s", webhook),
					"--mode=monitor",
					fmt.Sprintf("--name=%s", name),
					fmt.Sprintf("--commitUrl=%s", commitUrl),
				},
			},
		},
		Tags: []string{"slackbot"},
		// Add a timeout for the slackbot that is equal to the timeout for the current build
		Timeout: thisBuild.Timeout,
	}

	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	cr := svc.Projects.Builds.Create(project, b)
	_, err = cr.Do()
	if err != nil {
		log.Fatalf("Failed to create watcher build: %v", err)
	} else {
		log.Printf("Triggered watcher build.")
	}
}
