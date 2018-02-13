package packager

import "github.com/pkg/errors"

// Config holds our build configuration.
type Config struct {
	Project     string      `json:"project"     yaml:"project"`
	Push        bool        `json:"push"        yaml:"push"`
	Credentials Credentials `json:"credentials" yaml:"credentials"`
}

// Credentials contain our Docker Hub username and password.
type Credentials struct {
	Registry string `json:"registry" yaml:"registry"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// Validate ensures that all the build configuration was provided.
func (config Config) Validate() error {
	if config.Project == "" {
		return errors.New("project not found")
	}

	return config.Credentials.Validate()
}

// Validate ensures that all the credentials were provided.
func (credentials Credentials) Validate() error {
	if credentials.Registry == "" {
		return errors.New("registry not found")
	}

	return nil
}
