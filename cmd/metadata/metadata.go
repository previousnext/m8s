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

	// AnnotationCircleCIBuildNum
	AnnotationCircleCIBuildNum = "circleci.com/build/num"

	// AnnotationCircleCIBuildURL
	AnnotationCircleCIBuildURL = "circleci.com/build/url"

	// AnnotationCircleCICompareURL
	AnnotationCircleCICompareURL = "circleci.com/compare_url"

	// AnnotationCircleCISHA1
	AnnotationCircleCISHA1 = "circleci.com/sha1"

	// AnnotationCircleCIUsername
	AnnotationCircleCIUsername = "circleci.com/username"

	// AnnotationCircleCIBranch
	AnnotationCircleCIBranch = "circleci.com/branch"

	// AnnotationCircleCIJob
	AnnotationCircleCIJob = "circleci.com/job"

	// AnnotationCircleCIWorkflowID
	AnnotationCircleCIWorkflowID = "circleci.com/workflow/id"

	// AnnotationCircleCIWorkflowJobID
	AnnotationCircleCIWorkflowJobID = "circleci.com/workflow/job_id"

	// AnnotationCircleCIWorkflowWorkspaceID
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
