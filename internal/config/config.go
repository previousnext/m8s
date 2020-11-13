package config

import "io/ioutil"

// Load the Docker Compose file.
func Load(filePaths []string) (Config, error) {
	var dc Config

	for _, path := range filePaths {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return dc, err
		}

		err = yaml.Unmarshal(file, &dc)
		if err != nil {
			return dc, err
		}
	}

	return dc, nil
}

// Config is an object which encapsulates a Docker Compose file.
type Config struct {
	Services map[string]Service
}

// Service a service declared in a Docker Compose file.
type Service struct {
	Image       string            `yaml:"image"`
	Build       string            `yaml:"build"`
	Volumes     []string          `yaml:"volumes"`
	Entrypoint  []string          `yaml:"entrypoint"`
	Ports       []string          `yaml:"ports"`
	Environment []string          `yaml:"environment"`
	CapAdd      []string          `yaml:"cap_add"`
	Tmpfs       []string          `yaml:"tmpfs"`
	ExtraHosts  []string          `yaml:"extra_hosts"`
	Labels      map[string]string `yaml:"labels"`
}
