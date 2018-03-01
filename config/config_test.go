package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"time"
)

func TestLoad(t *testing.T) {
	have, err := Load("testdata/m8s.yml")
	assert.Nil(t, err)

	retention, err := time.ParseDuration("30m")
	assert.Nil(t, err)

	want := Config{
		Namespace: "test",
		Retention: retention,
		Auth: Auth{
			User: "nick",
			Pass: "rocks",
		},
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
