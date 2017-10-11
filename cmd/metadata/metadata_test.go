package metadata

import (
	"testing"

	pb "github.com/previousnext/m8s/pb"
	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	want := []*pb.Annotation{
		{
			Name:  AnnotationBitbucketBranch,
			Value: "master",
		},
		{
			Name:  AnnotationBitbucketRepoOwner,
			Value: "nick",
		},
		{
			Name:  AnnotationBitbucketRepoName,
			Value: "bar",
		},
		{
			Name:  AnnotationCircleCIRepositoryURL,
			Value: "http://example.com",
		},
		{
			Name:  AnnotationCircleCIPRNumber,
			Value: "1",
		},
		{
			Name:  AnnotationCircleCIPRUsername,
			Value: "nick",
		},
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
