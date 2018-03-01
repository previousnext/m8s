package environ

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	want := []string{
		"FOO=foo",
		"BAR=bar",
		"BAZ=baz",
	}

	have := filter([]string{
		"M8S_ENV_FOO=foo",
		"M8S_ENV_BAR=bar",
		"M8S_ENV_BAZ=baz",
		"WAH=wah",
		"WAZ=waz",
		"STUFF=stuff",
	})

	assert.Equal(t, want, have)
}
