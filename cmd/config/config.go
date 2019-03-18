package config

import (
	"io/ioutil"

	"github.com/previousnext/compose"
	"gopkg.in/yaml.v2"
)

// Config represents the values from the configuraiton file.
type Config struct {
	Init  []InitStep `yaml:"init"`
	Build []string   `yaml:"build"`
}

// InitStep represents a task that should be done in an init container.
type InitStep struct {
	Name      string                         `yaml:"name"`
	Image     string                         `yaml:"image"`
	Steps     []string                       `yaml:"steps"`
	Resources compose.ServiceDeployResources `yaml:"resources"`
	Volumes   []string                       `yaml:"volumes"`
}

// Load unmarshalls a Config object from a file.
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
