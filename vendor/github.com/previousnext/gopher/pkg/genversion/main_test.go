package genversion

import (
	"bytes"
	"io/ioutil"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestRenderVersionFile(t *testing.T) {
	var remoteFilepathTests = []struct {
		comment string
		want    string
		params  VersionConstants
	}{
		{
			"case 1",
			"./test-data/case-1.txt",
			VersionConstants{
				BuildDate:    "2018-02-02T12:00:00+11:00",
				BuildVersion: "v1.0.0-e43fd7",
			},
		},
		{
			"case 2",
			"./test-data/case-2.txt",
			VersionConstants{
				BuildDate:    "2020-11-11T12:00:00+11:00",
				BuildVersion: "v4.3.2",
			},
		},
	}

	for _, testCase := range remoteFilepathTests {
		want, err := ioutil.ReadFile(testCase.want)
		assert.Nil(t, err)

		tmpl, err := template.New("version.go.tmpl").Parse(templateContent)
		assert.Nil(t, err)

		actual := new(bytes.Buffer)
		RenderVersionFile(actual, tmpl, testCase.params)
		assert.Equal(t, string(want), actual.String(), testCase.comment)
	}
}
