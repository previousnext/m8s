package main

import (
	"fmt"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/api/core/v1"
)

var (
	cliGithubToken = kingpin.Flag("githubToken", "The Github access token.").Required().String()
	prCache        map[string][]*github.PullRequest
)

const (
	KeyCircleProjectRepoName = "circleci.com/project/reponame"
	KeyCircleProjectUserName = "circleci.com/project/username"
	KeyCircleBranch          = "circleci.com/branch"
)

func init() {
	prCache = make(map[string][]*github.PullRequest)
}

func main() {
	kingpin.Parse()
	// creates the in-cluster config.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *cliGithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientGithub := github.NewClient(tc)
	slayEmAll, err := getPodsToSlay(pods, clientGithub, ctx)
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range slayEmAll {
		fmt.Printf("Slaying Pod %s in namespace %s\n", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)
		err = clientset.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
}

// getPodsToSlay gathers any pods with closed PRs.
func getPodsToSlay(pods *v1.PodList, clientGithub *github.Client, ctx context.Context) ([]v1.Pod, error) {
	var podsToKill []v1.Pod
	for _, pod := range pods.Items {
		if _, ok := pod.Annotations[KeyCircleProjectRepoName]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, KeyCircleProjectRepoName)
			continue
		}
		if _, ok := pod.Annotations[KeyCircleProjectUserName]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, KeyCircleProjectUserName)
			continue
		}
		if _, ok := pod.Annotations[KeyCircleBranch]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, KeyCircleBranch)
			continue
		}

		prs, err := getClosedPrs(clientGithub, ctx, pod.Annotations[KeyCircleProjectUserName], pod.Annotations[KeyCircleProjectRepoName])
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			if pod.Annotations[KeyCircleBranch] == *pr.Head.Ref {
				podsToKill = append(podsToKill, pod)
			}
		}
	}

	return podsToKill, nil
}

// getClosedPrs gathers all closed PRs for an owner and repo.
func getClosedPrs(client *github.Client, ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	key := fmt.Sprintf("%s.%s", owner, repo)
	if val, ok := prCache[key]; ok {
		fmt.Printf("Loading PRs from cache for owner %s and repo %s\n", owner, repo)
		return val, nil
	}
	var allPrs []*github.PullRequest
	opt := &github.PullRequestListOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allPrs = append(allPrs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	prCache[key] = allPrs
	return allPrs, nil
}
