package main

import (
	"context"
	"os"
	"sync"

	"github.com/jrcasso/tugboat/providers/github"
	"github.com/jrcasso/tugboat/providers/kubernetes"
	"github.com/jrcasso/tugboat/tugboat"
	log "github.com/sirupsen/logrus"
)

const SERVICE_DIR = "services"

func main() {
	initializeLogging()
	var wg sync.WaitGroup
	services := tugboat.LoadServices(SERVICE_DIR)
	ctx := context.Background()
	plans := [][]tugboat.ExecutionPlan{}

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
		plan := provider.Plan(services)
		log.Infof("Generated execution plan: %v", plan)
		plans = append(plans, plan)
	}

	for _, plan := range plans {
		wg.Add(1)
		go tugboat.Execute(plan, &wg)
	}

	// Cleanup
	wg.Wait()
	for _, provider := range providers {
		for _, service := range services {
			wg.Add(1)
			go func(service tugboat.Service) {
				defer wg.Done()
				log.Infof("Cleaning up %v", service.Name)
				provider.Delete(service.Name)
			}(service)
		}
	}
	wg.Wait()
	log.Infof("Finished cleanup!")
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
