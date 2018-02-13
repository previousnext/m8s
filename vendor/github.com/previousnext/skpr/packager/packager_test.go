package packager

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"strings"
)

// Requires you to be running a local registry.
// See this projects docker-compose file.
func TestBuildMany(t *testing.T) {
	params := BuildManyParams{
		Config: Config{
			Project: "test",
			Credentials: Credentials{
				Registry: "localhost:5000",
				Username: "testuser",
				Password: "testuser",
			},
		},
		Images: []Image{
			{
				Name:       "frontend",
				Repository: "example/project",
				Tag:        "0.0.1-frontend",
				Dockerfile: "Dockerfile",
				Context:    "test/frontend",
			},
			{
				Name:       "backend",
				Repository: "example/project",
				Tag:        "0.0.1-backend",
				Dockerfile: "Dockerfile",
				Context:    "test/frontend",
			},
		},
	}

	// We will write all our output to this so we can inspect in our test suite.
	var buffer bytes.Buffer

	err := BuildMany(params, &buffer)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(buffer.String(), "Successfully built"))
}
