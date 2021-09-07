package main

import (
	"context"
	"os"

	"github.com/jrcasso/tugboat/providers/github"
	"github.com/jrcasso/tugboat/providers/kubernetes"
	"github.com/jrcasso/tugboat/tugboat"
	log "github.com/sirupsen/logrus"
)

const SERVICE_DIR = "services"

func main() {
	initializeLogging()
	services := tugboat.LoadServices(SERVICE_DIR)
	ctx := context.Background()

	githubClient := github.CreateClient(ctx)
	k8sClient := kubernetes.CreateClient()

	// Make this allocation dynamic based out of LoadServices
	providers := []tugboat.Provider{
		kubernetes.KubernetesProvider{
			Client:  k8sClient,
			Context: &ctx,
		},
		github.GithubProvider{
			Client:  &githubClient,
			Context: &ctx,
		},
	}

	for _, provider := range providers {
		// Add concurrency
		plan := provider.Plan(services)
		log.Infof("Execution plan: %v", plan)
		tugboat.Execute(plan)
	}

	// Cleanup
	for _, provider := range providers {
		for _, service := range services {
			provider.Delete(service.Name)
		}
	}
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
