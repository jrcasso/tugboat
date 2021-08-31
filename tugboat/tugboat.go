package tugboat

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Name      string `yaml:"name"`
	Namespace bool   `yaml:"Namespace"`
	Template  string `yaml:"template,omitempty"`
}

type Provider interface {
	Create(name string)
	// Retrieve()
	// Update()
	// Delete()

	// resolve()
}

func LoadServices(dir string) []Service {
	var configs []Service

	log.Debugf("Reading directory at: %+v", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read service directory: %+v", err)
	}

	log.Debugf("Found %+v files", len(files))
	for _, file := range files {
		config := processConfig(fmt.Sprintf("./%+v/%+v", dir, file.Name()))
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