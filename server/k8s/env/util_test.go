package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	assert.Equal(t, "cache", machine("/cache"))
}

func TestRetentionToUnix(t *testing.T) {
	now := time.Unix(1405544146, 0)

	unix, err := retentionToUnix(now, "24h")
	assert.Nil(t, err)

	assert.Equal(t, "1405630546", unix)
}
