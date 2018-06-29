package pod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCodePath(t *testing.T) {
	exists, path := FindCodePath([]string{
		".:/data",
	})
	assert.True(t, exists)
	assert.Equal(t, "/data", path)
}
