package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheDirectories(t *testing.T) {
	want := []CacheDirectory{
		{
			Name: "composer",
			Path: "/root/.composer",
		},
		{
			Name: "yarn",
			Path: "/usr/local/share/.cache/yarn",
		},
	}

	caches, err := cacheDirectories("composer:/root/.composer,yarn:/usr/local/share/.cache/yarn")
	assert.Nil(t, err)

	assert.Equal(t, want, caches)
}
