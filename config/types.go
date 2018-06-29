package config

import "time"

// Config which is passed to the client.
type Config struct {
	Namespace string        `yaml:"namespace" json:"namespace"`
	Retention time.Duration `yaml:"retention" json:"retention"`
	Auth      Auth          `yaml:"auth"      json:"auth"`
	Build     Build         `yaml:"build"     json:"build"`
	Cache     Cache         `yaml:"cache"     json:"cache"`
	Secrets   Secrets       `yaml:"secrets"   json:"secrets"`
	Port      int           `yaml:"port"      json:"port"`
}

// Auth to secure the environment.
type Auth struct {
	User string `yaml:"user" json:"user"`
	Pass string `yaml:"pass" json:"pass"`
}

// Build steps to be run after the environment is up.
type Build struct {
	Container string   `yaml:"container" json:"container"`
	Steps     []string `yaml:"steps"     json:"steps"`
}

// Cache directories for between builds.
type Cache struct {
	Type  string   `yaml:"type"  json:"type"`
	Paths []string `yaml:"paths" json:"paths"`
}

// Secrets for interacting with private repositories.
type Secrets struct {
	DockerCfg string `yaml:"dockercfg" json:"dockercfg"`
	SSH       string `yaml:"ssh"       json:"ssh"`
}
