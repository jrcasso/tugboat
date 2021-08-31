package main

import (
	"context"
	"os"
	"time"

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

	providers := []tugboat.Provider{
		github.GithubProvider{
			Client:  &githubClient,
			Context: &ctx,
		},
	}

	for _, provider := range providers {
		for _, service := range services {
			provider.Create(service.Name)
			github.GetOrgRepositories(ctx, githubClient)
			time.Sleep(3)
			github.DeleteRepository(ctx, githubClient, service.Name)
		}
	}

	namespaces := kubernetes.GetNamespaces(*k8sClient)
	for _, namespace := range namespaces.Items {
		log.Debugf("Found namespace: %+v", namespace.Name)
	}
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
