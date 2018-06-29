package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	have, err := Load("testdata/m8s.yml")
	assert.Nil(t, err)

	retention, err := time.ParseDuration("30m")
	assert.Nil(t, err)

	want := Config{
		Namespace: "test",
		Retention: retention,
		Build: Build{
			Container: "app",
			Steps: []string{
				"echo 1",
				"echo 2",
			},
		},
		Cache: Cache{
			Type: "standard",
			Paths: []string{
				"/root/.composer",
			},
		},
		Secrets: Secrets{
			DockerCfg: "test-dockercfg",
			SSH:       "test-ssh",
		},
	}

	assert.Equal(t, want, have)
}
