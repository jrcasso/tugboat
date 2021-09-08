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
		wg.Add(1)
		go func(provider tugboat.Provider, service []tugboat.Service) {
			plan := provider.Plan(services, &wg)
			log.Infof("Generated execution plan: %+v", plan)
			plans = append(plans, plan)
		}(provider, services)
	}
	log.Debugf("Waiting for execution planning to complete...")
	wg.Wait()
	log.Debugf("Planning completed!")

	for _, plan := range plans {
		wg.Add(1)
		go tugboat.Execute(plan, &wg)
	}
	log.Debugf("Waiting for plan execution to complete...")
	wg.Wait()
	log.Debugf("Plan execution completed!")

	// Cleanup
	for _, provider := range providers {
		for _, service := range services {
			wg.Add(1)
			go func(provider tugboat.Provider, service tugboat.Service) {
				defer wg.Done()
				log.Infof("Cleaning up %v", service.Name)
				provider.Delete(service.Name)
			}(provider, service)
		}
	}
	log.Debugf("Waiting for cleanup to complete...")
	wg.Wait()
	log.Debugf("Cleanup completed!")
}

func initializeLogging() {
	// TODO: Implement dynamic log level, output, and format switches
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}
