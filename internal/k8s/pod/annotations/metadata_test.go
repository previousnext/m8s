package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	want := map[string]string{
		AnnotationBitbucketBranch: "master",
		AnnotationBitbucketRepoOwner: "nick",
		AnnotationBitbucketRepoName: "bar",
		AnnotationBitbucketBuildNumber: "quiix",
		AnnotationBitbucketCommit: "quiix",
		AnnotationBitbucketRepoOwnerUUID: "quiix",
		AnnotationBitbucketRepoUUID: "quiix",
		AnnotationBitbucketTag: "quiix",
		AnnotationBitbucketBookmark: "quiix",
		AnnotationBitbucketParallelStep: "quiix",
		AnnotationBitbucketParallelStepCount: "quiix",
		AnnotationCircleCIRepositoryURL: "http://example.com",
		AnnotationCircleCIPRNumber: "1",
		AnnotationCircleCIPRUsername: "nick",
		AnnotationCircleCIBuildNum: "xxx",
		AnnotationCircleCIBuildURL: "xxx",
		AnnotationCircleCICompareURL: "xxx",
		AnnotationCircleCISHA1: "xxx",
		AnnotationCircleCIRepositoryName: "xxx",
		AnnotationCircleCIRepositoryUsername: "xxx",
		AnnotationCircleCIUsername: "xxx",
		AnnotationCircleCIBranch: "xxx",
		AnnotationCircleCIJob: "xxx",
		AnnotationCircleCIWorkflowID: "xxx",
		AnnotationCircleCIWorkflowJobID: "xxx",
		AnnotationCircleCIWorkflowWorkspaceID: "xxx",
	}

	have, err := Annotations([]string{
		"BITBUCKET_BRANCH=master",
		"BITBUCKET_REPO_OWNER=nick",
		"BITBUCKET_REPO_SLUG=bar",
		"BITBUCKET_BUILD_NUMBER=quiix",
		"BITBUCKET_COMMIT=quiix",
		"BITBUCKET_REPO_OWNER_UUID=quiix",
		"BITBUCKET_REPO_UUID=quiix",
		"BITBUCKET_TAG=quiix",
		"BITBUCKET_BOOKMARK=quiix",
		"BITBUCKET_PARALLEL_STEP=quiix",
		"BITBUCKET_PARALLEL_STEP_COUNT=quiix",
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
