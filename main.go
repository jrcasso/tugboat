package main

import (
	"context"
	"os"

	"github.com/jrcasso/tugboat/providers/github"

	log "github.com/sirupsen/logrus"
)

func main() {
	initializeLogging()
	ctx := context.Background()
	client := github.CreateClient(ctx)
	github.CreateRepository(ctx, client, "test-repo")
	github.DeleteRepository(ctx, client, "test-repo")
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
