package packager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	image := Image{}

	err := image.Validate()
	assert.Equal(t, err.Error(), "name not found")

	image.Name = "example"
	err = image.Validate()
	assert.Equal(t, err.Error(), "repository not found")

	image.Repository = "previousnext/example"
	err = image.Validate()
	assert.Equal(t, err.Error(), "image tag not found")

	image.Tag = "0.0.1"
	err = image.Validate()
	assert.Equal(t, err.Error(), "Dockerfile not found")

	image.Dockerfile = "Dockerfile"
	err = image.Validate()
	assert.Equal(t, err.Error(), "context not found")

	image.Context = "."
	err = image.Validate()
	assert.Nil(t, err)
}
