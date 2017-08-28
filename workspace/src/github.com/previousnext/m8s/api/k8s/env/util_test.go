package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	assert.Equal(t, "cache", machine("/cache"))
}
