package packager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	config := Config{}

	err := config.Validate()
	assert.Equal(t, err.Error(), "project not found")

	config.Project = "example"
	err = config.Validate()
	assert.Equal(t, err.Error(), "registry not found")

	config.Credentials.Registry = "example.com"
	err = config.Validate()
	assert.Nil(t, err)
}
