package metadata

import (
	"strings"

	pb "github.com/previousnext/m8s/pb"
)

const (
	// AnnotationBitbucketBranch is an identifier for Bitbucket branch this environment was built from.
	AnnotationBitbucketBranch = "bitbucket.org/branch"

	// AnnotationBitbucketRepoOwner is an identifier for Bitbucket repository owner.
	AnnotationBitbucketRepoOwner = "bitbucket.org/repo/owner"

	// AnnotationBitbucketRepoName is an identifier for Bitbucket repository this environment was built from.
	AnnotationBitbucketRepoName = "bitbucket.org/repo/name"

	// AnnotationBitbucketBuildNumber is the unique identifier for a build. It increments with each build and can be
	// used to create unique artifact names.
	AnnotationBitbucketBuildNumber = "bitbucket.org/build_number"

	// AnnotationBitbucketCommit is the commit hash of a commit that kicked off the build.
	AnnotationBitbucketCommit = "bitbucket.org/commit"

	// AnnotationBitbucketRepoOwnerUUID is the UUID of the account in which the repository lives.
	AnnotationBitbucketRepoOwnerUUID = "bitbucket.org/repo_owner_uuid"

	// AnnotationBitbucketRepoUUID is the UUID of the repository.
	AnnotationBitbucketRepoUUID = "bitbucket.org/repo_uuid"

	// AnnotationBitbucketTag is the tag of a commit that kicked off the build. This value is only available on tags.
	AnnotationBitbucketTag = "bitbucket.org/tag"

	// AnnotationBitbucketBookmark is for use with Mercurial projects.
	AnnotationBitbucketBookmark = "bitbucket.org/bookmark"

	// AnnotationBitbucketParallelStep is zero-based index of the current step in the group.
	AnnotationBitbucketParallelStep = "bitbucket.org/parallel_step"

	// AnnotationBitbucketParallelStepCount is Total number of steps in the group.
	AnnotationBitbucketParallelStepCount = "bitbucket.org/parallel_step_count"

	// AnnotationCircleCIRepositoryUsername is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryUsername = "circleci.com/project/username"

	// AnnotationCircleCIRepositoryName is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryName = "circleci.com/project/reponame"

	// AnnotationCircleCIRepositoryURL is an identifier for CircleCI repository the environment was built from.
	AnnotationCircleCIRepositoryURL = "circleci.com/repository/url"

	// AnnotationCircleCIPRNumber is an identifier for CircleCI pull request the environment was built from.
	AnnotationCircleCIPRNumber = "circleci.com/pr/number"

	// AnnotationCircleCIPRUsername is an identifier for CircleCI pull request submitted by a user.
	AnnotationCircleCIPRUsername = "circleci.com/pr/username"

	// AnnotationCircleCIBuildNum is the number of the CircleCI build.
	AnnotationCircleCIBuildNum = "circleci.com/build/num"

	// AnnotationCircleCIBuildURL is the URL for the current build.
	AnnotationCircleCIBuildURL = "circleci.com/build/url"

	// AnnotationCircleCICompareURL is the GitHub or Bitbucket URL to compare commits of a build.
	AnnotationCircleCICompareURL = "circleci.com/compare_url"

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
	AnnotationCircleCIWorkflowID = "circleci.com/workflow/id"

	// AnnotationCircleCIWorkflowJobID is the workflow job ID.
	AnnotationCircleCIWorkflowJobID = "circleci.com/workflow/job_id"

	// AnnotationCircleCIWorkflowWorkspaceID is the workspace ID.
	AnnotationCircleCIWorkflowWorkspaceID = "circleci.com/workflow/workspace_id"
)

// Annotations are used for attaching metadata to a environment.
func Annotations(envs []string) ([]*pb.Annotation, error) {
	var annotations []*pb.Annotation

	for _, env := range envs {
		sl := strings.Split(env, "=")

		if len(sl) != 2 {
			continue
		}

		switch sl[0] {
		// Check if we have Bitbucket Pipelines environment variables.
		// https://confluence.atlassian.com/bitbucket/environment-variables-794502608.html
		case "BITBUCKET_BRANCH":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketBranch, Value: sl[1]})
		case "BITBUCKET_REPO_OWNER":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketRepoOwner, Value: sl[1]})
		case "BITBUCKET_REPO_SLUG":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketRepoName, Value: sl[1]})
		case "BITBUCKET_BUILD_NUMBER":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketBuildNumber, Value: sl[1]})
		case "BITBUCKET_COMMIT":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketCommit, Value: sl[1]})
		case "BITBUCKET_REPO_OWNER_UUID":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketRepoOwnerUUID, Value: sl[1]})
		case "BITBUCKET_REPO_UUID":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketRepoUUID, Value: sl[1]})
		case "BITBUCKET_TAG":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketTag, Value: sl[1]})
		case "BITBUCKET_BOOKMARK":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketBookmark, Value: sl[1]})
		case "BITBUCKET_PARALLEL_STEP":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketParallelStep, Value: sl[1]})
		case "BITBUCKET_PARALLEL_STEP_COUNT":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationBitbucketParallelStepCount, Value: sl[1]})

		// Check if we have CircleCI environment variables.
		// https://circleci.com/docs/2.0/env-vars/
		case "CIRCLE_PR_NUMBER":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIPRNumber, Value: sl[1]})
		case "CIRCLE_PR_USERNAME":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIPRUsername, Value: sl[1]})
		case "CIRCLE_BUILD_NUM":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIBuildNum, Value: sl[1]})
		case "CIRCLE_BUILD_URL":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIBuildURL, Value: sl[1]})
		case "CIRCLE_COMPARE_URL":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCICompareURL, Value: sl[1]})
		case "CIRCLE_SHA1":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCISHA1, Value: sl[1]})
		case "CIRCLE_PROJECT_REPONAME":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIRepositoryName, Value: sl[1]})
		case "CIRCLE_PROJECT_USERNAME":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIRepositoryUsername, Value: sl[1]})
		case "CIRCLE_USERNAME":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIUsername, Value: sl[1]})
		case "CIRCLE_BRANCH":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIBranch, Value: sl[1]})
		case "CIRCLE_JOB":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIJob, Value: sl[1]})
		case "CIRCLE_REPOSITORY_URL":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIRepositoryURL, Value: sl[1]})
		case "CIRCLE_WORKFLOW_ID":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIWorkflowID, Value: sl[1]})
		case "CIRCLE_WORKFLOW_JOB_ID":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIWorkflowJobID, Value: sl[1]})
		case "CIRCLE_WORKFLOW_WORKSPACE_ID":
			annotations = append(annotations, &pb.Annotation{Name: AnnotationCircleCIWorkflowWorkspaceID, Value: sl[1]})
		}
	}

	return annotations, nil
}
