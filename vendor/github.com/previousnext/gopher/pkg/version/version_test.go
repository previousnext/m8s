package version

import (
	"testing"
	"bytes"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	var buffer bytes.Buffer

	params := PrintParams{
		OS: "leenux",
		Arch: "nintendo64",
	}

	err := Print(&buffer, params)
	assert.Equal(t, err.Error(), "version not found")

	params.Version = "0.0.1"
	err = Print(&buffer, params)
	assert.Equal(t, err.Error(), "commit not found")

	params.Commit = "abcdefg"
	err = Print(&buffer, params)
	assert.Nil(t, err)

	assert.True(t, strings.Contains(buffer.String(), params.Version))
	assert.True(t, strings.Contains(buffer.String(), params.Commit))
	assert.True(t, strings.Contains(buffer.String(), params.OS))
	assert.True(t, strings.Contains(buffer.String(), params.Arch))
}