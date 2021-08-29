package main

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var requiredEnvs = [...]string{
	"GITHUB_TOKEN",
	"GITHUB_ORGANIZATION",
}

func main() {
	initializeLogging()
	validateEnvironment()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:        github.String("test-repo"),
		Private:     github.Bool(true),
		Description: github.String("Test repository"),
	}
	repo, _, err := client.Repositories.Create(ctx, os.Getenv("GITHUB_ORGANIZATION"), repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully created new repo: %v\n", repo.GetName())
}

func validateEnvironment() {
	for i := 0; i < len(requiredEnvs); i++ {
		var env = os.Getenv(requiredEnvs[i])
		if env == "" {
			log.Fatalf("Required environment variable not set: %v", requiredEnvs[i])
		}
	}
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
