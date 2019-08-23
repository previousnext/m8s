package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	file := "./test-data/test.yaml"
	compose, err := Load(file)
	assert.Nil(t, err)

	assert.Equal(t, compose.Services["a"].Image, "hostname.io/org/repo:version")
	assert.Equal(t, compose.Services["a"].Volumes, []string{"a:/a", "b:/b"})
	assert.Equal(t, compose.Services["a"].Entrypoint, []string{"/entrypoint"})
	assert.Equal(t, compose.Services["a"].Ports, []string{"1000:2000"})
	assert.Equal(t, compose.Services["a"].Environment, []string{"SOME_VAR=somevalue"})
	assert.Equal(t, compose.Services["a"].CapAdd, []string{"NET_ADMIN"})
	assert.Equal(t, compose.Services["a"].Tmpfs, []string{"/tmp"})
	assert.Equal(t, compose.Services["a"].Deploy.Resources.Reservations.CPUs, "50m")
	assert.Equal(t, compose.Services["a"].Deploy.Resources.Reservations.Memory, "768Mi")
	assert.Equal(t, compose.Services["a"].Deploy.Resources.Limits.CPUs, "500m")
	assert.Equal(t, compose.Services["a"].Deploy.Resources.Limits.Memory, "2048Mi")
	assert.Equal(t, compose.Services["a"].ExtraHosts, []string{"some.hostname:1.2.3.4"})
	assert.Equal(t, compose.Services["a"].Labels, map[string]string{ "m8s.io/skip": "true"})
}
