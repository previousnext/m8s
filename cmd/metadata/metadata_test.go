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
		{
			Name: AnnotationCircleCIBuildNum,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIBuildURL,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCICompareURL,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCISHA1,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIRepositoryName,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIRepositoryUsername,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIUsername,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIBranch,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIJob,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIWorkflowID,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIWorkflowJobID,
			Value: "xxx",
		},
		{
			Name: AnnotationCircleCIWorkflowWorkspaceID,
			Value: "xxx",
		},
	}

	have, err := Annotations([]string{
		"BITBUCKET_BRANCH=master",
		"BITBUCKET_REPO_OWNER=nick",
		"BITBUCKET_REPO_SLUG=bar",
		"CIRCLE_REPOSITORY_URL=http://example.com",
		"CIRCLE_PR_NUMBER=1",
		"CIRCLE_PR_USERNAME=nick",
		"CIRCLE_BUILD_NUM=xxx",
		"CIRCLE_BUILD_URL=xxx",
		"CIRCLE_COMPARE_URL=xxx",
		"CIRCLE_SHA1=xxx",
		"CIRCLE_PROJECT_REPONAME=xxx",
		"CIRCLE_PROJECT_USERNAME=xxx",
		"CIRCLE_USERNAME=xxx",
		"CIRCLE_BRANCH=xxx",
		"CIRCLE_JOB=xxx",
		"CIRCLE_WORKFLOW_ID=xxx",
		"CIRCLE_WORKFLOW_JOB_ID=xxx",
		"CIRCLE_WORKFLOW_WORKSPACE_ID=xxx",
		"FOO=bar",
	})
	assert.Nil(t, err)

	assert.Equal(t, want, have)
}
