package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jrcasso/tugboat/providers/github"
	"github.com/jrcasso/tugboat/providers/kubernetes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const SERVICE_DIR = "services"

type Service struct {
	Name      string `yaml:"name"`
	Namespace bool   `yaml:"Namespace"`
	Template  string `yaml:"template,omitempty"`
}

func main() {
	initializeLogging()
	services := loadServiceConfigs()
	ctx := context.Background()

	githubClient := github.CreateClient(ctx)
	k8sClient := kubernetes.CreateClient()

	for _, service := range services {
		github.CreateRepository(ctx, githubClient, service.Name)
		github.GetOrgRepositories(ctx, githubClient)
		time.Sleep(1)
		github.DeleteRepository(ctx, githubClient, service.Name)
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

func loadServiceConfigs() []Service {
	var configs []Service

	log.Debugf("Reading directory at: %+v", SERVICE_DIR)
	files, err := os.ReadDir(SERVICE_DIR)
	if err != nil {
		log.Fatalf("Failed to read service directory: %+v", err)
	}

	log.Debugf("Found %+v files", len(files))
	for _, file := range files {
		config := processConfig(fmt.Sprintf("./%+v/%+v", SERVICE_DIR, file.Name()))
		log.Debugf("Found configuration: %+v", config)
		configs = append(configs, config)
	}

	return configs
}

func processConfig(path string) Service {
	t := Service{}

	log.Debugf("Reading configuration at: %+v", path)
	configBytes, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Fatalf("Failed to read service configuration at %+v: %+v", path, readErr)
	}

	log.Debugf("Unmarshalling configuration at: %+v", path)
	yamlErr := yaml.Unmarshal(configBytes, &t)
	if yamlErr != nil {
		log.Fatalf("Failed to unmarshal service configuration at %+v: %+v", path, yamlErr)
	}

	return t
}
