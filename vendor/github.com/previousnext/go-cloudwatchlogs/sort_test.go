package cloudwatchlogs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeLogs(t *testing.T) {
	first := &Log{
		Stream:    "one",
		Timestamp: time.Unix(1451098534, 0),
		Message:   "This is the first record",
	}
	second := &Log{
		Stream:    "two",
		Timestamp: time.Unix(1451098535, 0),
		Message:   "This is the second record",
	}
	third := &Log{
		Stream:    "one",
		Timestamp: time.Unix(1451098536, 0),
		Message:   "This is the third record",
	}

	// This is our initial set of records. Soon to be filtered against a new set.
	set1 := Logs{
		first,
		third,
	}

	// This is a set which will slot in between the 2 first logs.
	set2 := Logs{
		second,
	}

	// Merge and compare.
	set3 := Logs{
		first,
		second,
		third,
	}
	assert.Equal(t, set3, MergeLogs(set1, set2), "Successfully ordered our logs.")
}
