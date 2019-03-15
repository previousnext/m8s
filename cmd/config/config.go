package config

import (
	"github.com/previousnext/compose"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Init  []InitStep `yaml:"init"`
	Build []string   `yaml:"build"`
}

type InitStep struct {
	Name      string                         `yaml:"name"`
	Image     string                         `yaml:"image"`
	Steps     []string                       `yaml:"steps"`
	Resources compose.ServiceDeployResources `yaml:"resources"`
	Volumes   []string                       `yaml:"volumes"`
}

// Helper function to load testing steps.
func Load(f string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
