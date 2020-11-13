package annotations

import (
	"strings"
)

const (
	// AnnotationBitbucketBranch is an identifier for Bitbucket branch this environment was built from.
	AnnotationBitbucketBranch = "bitbucket.org/branch"

	// AnnotationBitbucketRepoOwner is an identifier for Bitbucket repository owner.
	AnnotationBitbucketRepoOwner = "bitbucket.org/repo-owner"

	// AnnotationBitbucketRepoName is an identifier for Bitbucket repository this environment was built from.
	AnnotationBitbucketRepoName = "bitbucket.org/repo-name"

	// AnnotationBitbucketBuildNumber is the unique identifier for a build. It increments with each build and can be
	// used to create unique artifact names.
	AnnotationBitbucketBuildNumber = "bitbucket.org/build-number"

	// AnnotationBitbucketCommit is the commit hash of a commit that kicked off the build.
	AnnotationBitbucketCommit = "bitbucket.org/commit"

	// AnnotationBitbucketRepoOwnerUUID is the UUID of the account in which the repository lives.
	AnnotationBitbucketRepoOwnerUUID = "bitbucket.org/repo-owner-uuid"

	// AnnotationBitbucketRepoUUID is the UUID of the repository.
	AnnotationBitbucketRepoUUID = "bitbucket.org/repo-uuid"

	// AnnotationBitbucketTag is the tag of a commit that kicked off the build. This value is only available on tags.
	AnnotationBitbucketTag = "bitbucket.org/tag"

	// AnnotationBitbucketBookmark is for use with Mercurial projects.
	AnnotationBitbucketBookmark = "bitbucket.org/bookmark"

	// AnnotationBitbucketParallelStep is zero-based index of the current step in the group.
	AnnotationBitbucketParallelStep = "bitbucket.org/parallel-step"

	// AnnotationBitbucketParallelStepCount is Total number of steps in the group.
	AnnotationBitbucketParallelStepCount = "bitbucket.org/parallel-step-count"

	// AnnotationCircleCIRepositoryUsername is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryUsername = "circleci.com/project-username"

	// AnnotationCircleCIRepositoryName is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryName = "circleci.com/project-reponame"

	// AnnotationCircleCIRepositoryURL is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryURL = "circleci.com/repository-url"

	// AnnotationCircleCIPRNumber is an identifier for CircleCI pull request the environment was built from.
	AnnotationCircleCIPRNumber = "circleci.com/pr-number"

	// AnnotationCircleCIPRUsername is an identifier for CircleCI pull request submitted by a user.
	AnnotationCircleCIPRUsername = "circleci.com/pr-username"

	// AnnotationCircleCIBuildNum is the number of the CircleCI build.
	AnnotationCircleCIBuildNum = "circleci.com/build-num"

	// AnnotationCircleCIBuildURL is the URL for the current build.
	AnnotationCircleCIBuildURL = "circleci.com/build-url"

	// AnnotationCircleCICompareURL is the GitHub or Bitbucket URL to compare commits of a build.
	AnnotationCircleCICompareURL = "circleci.com/compare-url"

	// AnnotationCircleCISHA1 is the SHA1 hash of the last commit of the current build.
	AnnotationCircleCISHA1 = "circleci.com/sha1"

	// AnnotationCircleCIUsername is the GitHub or Bitbucket username of the user who triggered the build.
	AnnotationCircleCIUsername = "circleci.com/username"

	// AnnotationCircleCIBranch is the name of the Git branch currently being built.
	AnnotationCircleCIBranch = "circleci.com/branch"

	// AnnotationCircleCIJob is the name of the current job.
	AnnotationCircleCIJob = "circleci.com/job"

	// AnnotationCircleCIWorkflowID is a unique identifier for the workflow instance of the current job. This
	// identifier is the same for every job in a given workflow instance.
	AnnotationCircleCIWorkflowID = "circleci.com/workflow-id"

	// AnnotationCircleCIWorkflowJobID is the workflow job ID.
	AnnotationCircleCIWorkflowJobID = "circleci.com/workflow-job-id"

	// AnnotationCircleCIWorkflowWorkspaceID is the workspace ID.
	AnnotationCircleCIWorkflowWorkspaceID = "circleci.com/workflow-workspace-id"
)

// FromEnvironment will get metadata from the environment which can be used as annotations.
func FromEnvironment(envs []string) map[string]string {
	annotations := make(map[string]string)

	for _, env := range envs {
		sl := strings.Split(env, "=")

		if len(sl) != 2 {
			continue
		}

		switch sl[0] {
		// Check if we have Bitbucket Pipelines environment variables.
		// https://confluence.atlassian.com/bitbucket/environment-variables-794502608.html
		case "BITBUCKET_BRANCH":
			annotations[AnnotationBitbucketBranch] = sl[1]
		case "BITBUCKET_REPO_OWNER":
			annotations[AnnotationBitbucketRepoOwner] = sl[1]
		case "BITBUCKET_REPO_SLUG":
			annotations[AnnotationBitbucketRepoName] = sl[1]
		case "BITBUCKET_BUILD_NUMBER":
			annotations[AnnotationBitbucketBuildNumber] = sl[1]
		case "BITBUCKET_COMMIT":
			annotations[AnnotationBitbucketCommit] = sl[1]
		case "BITBUCKET_REPO_OWNER_UUID":
			annotations[AnnotationBitbucketRepoOwnerUUID] = sl[1]
		case "BITBUCKET_REPO_UUID":
			annotations[AnnotationBitbucketRepoUUID] = sl[1]
		case "BITBUCKET_TAG":
			annotations[AnnotationBitbucketTag] = sl[1]
		case "BITBUCKET_BOOKMARK":
			annotations[AnnotationBitbucketBookmark] = sl[1]
		case "BITBUCKET_PARALLEL_STEP":
			annotations[AnnotationBitbucketParallelStep] = sl[1]
		case "BITBUCKET_PARALLEL_STEP_COUNT":
			annotations[AnnotationBitbucketParallelStepCount] = sl[1]

		// Check if we have CircleCI environment variables.
		// https://circleci.com/docs/2.0/env-vars/
		case "CIRCLE_PR_NUMBER":
			annotations[AnnotationCircleCIPRNumber] = sl[1]
		case "CIRCLE_PR_USERNAME":
			annotations[AnnotationCircleCIPRUsername] = sl[1]
		case "CIRCLE_BUILD_NUM":
			annotations[AnnotationCircleCIBuildNum] = sl[1]
		case "CIRCLE_BUILD_URL":
			annotations[AnnotationCircleCIBuildURL] = sl[1]
		case "CIRCLE_COMPARE_URL":
			annotations[AnnotationCircleCICompareURL] = sl[1]
		case "CIRCLE_SHA1":
			annotations[AnnotationCircleCISHA1] = sl[1]
		case "CIRCLE_PROJECT_REPONAME":
			annotations[AnnotationCircleCIRepositoryName] = sl[1]
		case "CIRCLE_PROJECT_USERNAME":
			annotations[AnnotationCircleCIRepositoryUsername] = sl[1]
		case "CIRCLE_USERNAME":
			annotations[AnnotationCircleCIUsername] = sl[1]
		case "CIRCLE_BRANCH":
			annotations[AnnotationCircleCIBranch] = sl[1]
		case "CIRCLE_JOB":
			annotations[AnnotationCircleCIJob] = sl[1]
		case "CIRCLE_REPOSITORY_URL":
			annotations[AnnotationCircleCIRepositoryURL] = sl[1]
		case "CIRCLE_WORKFLOW_ID":
			annotations[AnnotationCircleCIWorkflowID] = sl[1]
		case "CIRCLE_WORKFLOW_JOB_ID":
			annotations[AnnotationCircleCIWorkflowJobID] = sl[1]
		case "CIRCLE_WORKFLOW_WORKSPACE_ID":
			annotations[AnnotationCircleCIWorkflowWorkspaceID] = sl[1]
		}
	}

	return annotations
}
