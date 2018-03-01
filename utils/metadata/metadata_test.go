package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	want := map[string]string{
		AnnotationBitbucketBranch:       "master",
		AnnotationBitbucketRepoOwner:    "nick",
		AnnotationBitbucketRepoName:     "bar",
		AnnotationCircleCIRepositoryURL: "http://example.com",
		AnnotationCircleCIPRNumber:      "1",
		AnnotationCircleCIPRUsername:    "nick",
	}

	have, err := Annotations([]string{
		"BITBUCKET_BRANCH=master",
		"BITBUCKET_REPO_OWNER=nick",
		"BITBUCKET_REPO_SLUG=bar",
		"CIRCLE_REPOSITORY_URL=http://example.com",
		"CIRCLE_PR_NUMBER=1",
		"CIRCLE_PR_USERNAME=nick",
		"FOO=bar",
	})
	assert.Nil(t, err)

	assert.Equal(t, want, have)
}
