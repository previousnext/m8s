package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	assert.Equal(t, "cache", Machine("/cache"))
}
