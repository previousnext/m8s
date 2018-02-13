package packager

import "github.com/pkg/errors"

// Image declares a Dockerfile and Context which needs to be built.
type Image struct {
	Name       string `json:"name"       yaml:"name"`
	Repository string `json:"repository" yaml:"repository"`
	Tag        string `json:"tag"        yaml:"tag"`
	Dockerfile string `json:"dockerfile" yaml:"dockerfile"`
	Context    string `json:"context"    yaml:"context"`
}

// Validate ensures that all the image configuration was provided.
func (image Image) Validate() error {
	if image.Name == "" {
		return errors.New("name not found")
	}

	if image.Repository == "" {
		return errors.New("repository not found")
	}

	if image.Tag == "" {
		return errors.New("image tag not found")
	}

	if image.Dockerfile == "" {
		return errors.New("Dockerfile not found")
	}

	if image.Context == "" {
		return errors.New("context not found")
	}

	return nil
}
