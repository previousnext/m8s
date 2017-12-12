package compose

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// DockerCompose is an object which encapsulates a Docker Compose file.
type DockerCompose struct {
	Services map[string]Service
}

// Service a service declared in a Docker Compose file.
type Service struct {
	Image       string        `yaml:"image"`
	Build       string        `yaml:"build"`
	Volumes     []string      `yaml:"volumes"`
	Entrypoint  []string      `yaml:"entrypoint"`
	Ports       []string      `yaml:"ports"`
	Environment []string      `yaml:"environment"`
	CapAdd      []string      `yaml:"cap_add"`
	Tmpfs       []string      `yaml:"tmpfs"`
	Deploy      ServiceDeploy `yaml:"deploy"`
}

// ServiceDeploy provides deployment information for a service.
type ServiceDeploy struct {
	Resources ServiceDeployResources `yaml:"resources"`
}

// ServiceDeployResources provides deployment resources information for a service.
type ServiceDeployResources struct {
	Limits       ServiceDeployResource `yaml:"limits"`
	Reservations ServiceDeployResource `yaml:"reservations"`
}

// ServiceDeployResource provides a single deployment resource information for a service.
type ServiceDeployResource struct {
	CPUs   string `yaml:"cpus"`
	Memory string `yaml:"memory"`
}

// Load the Docker Compose file.
func Load(path string) (DockerCompose, error) {
	var dc DockerCompose

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return dc, err
	}

	err = yaml.Unmarshal(file, &dc)
	if err != nil {
		return dc, err
	}

	return dc, nil
}
