package tugboat

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ExecutionPlan struct {
	Function  func(string)
	Arguments string
}

type Service struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Repo      string `yaml:"repo"`
	Template  string `yaml:"template,omitempty"`
}

type Provider interface {
	Create(name string)
	Retrieve() []string
	// Update()
	Delete(name string)
	Execute(plan []ExecutionPlan)
	Plan(services []Service) []ExecutionPlan
}

func LoadServices(dir string) []Service {
	var configs []Service

	log.Debugf("Reading directory at: %+v", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read service directory: %+v", err)
	}

	log.Debugf("Found %+v files", len(files))
	dnsRegex, _ := regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	for _, file := range files {
		config := processConfig(fmt.Sprintf("./%+v/%+v", dir, file.Name()))
		log.Debugf("Found configuration: %+v", config)

		matched := dnsRegex.MatchString(config.Name)
		if err != nil {
			log.Fatalf("Failed to enforce regex: %+v", err)
		}
		if !matched {
			log.Errorf("Skipping configuration with non-compliant DNS name: %+v", config.Name)
			continue
		}

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
