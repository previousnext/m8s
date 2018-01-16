package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderVersionOutput(t *testing.T) {
	var remoteFilepathTests = []struct {
		comment string
		want    string
		params  cmdVersion
	}{
		{
			"case 1",
			`Version  v1.0.0
Date     2018-01-01T00:00:00+11:00
API      v1
OS       darwin
Arch     amd64`,
			cmdVersion{
				APICompatibility: 1,
				BuildDate:        "2018-01-01T00:00:00+11:00",
				BuildVersion:     "v1.0.0",
				GOARCH:           "amd64",
				GOOS:             "darwin",
			},
		},
		{
			"case 2",
			`Version  v1.0.0-b715353
Date     2018-01-12T16:42:06+11:00
API      v2
OS       linux
Arch     arm6`,
			cmdVersion{
				APICompatibility: 2,
				BuildDate:        "2018-01-12T16:42:06+11:00",
				BuildVersion:     "v1.0.0-b715353",
				GOARCH:           "arm6",
				GOOS:             "linux",
			},
		},
	}

	for _, testCase := range remoteFilepathTests {
		actual := renderVersionOutput(&testCase.params)
		assert.Equal(t, testCase.want, actual, testCase.comment)
	}
}
