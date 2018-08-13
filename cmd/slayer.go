package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/previousnext/m8s/cmd/metadata"
)

var (
	slayCache map[string][]*github.PullRequest
)

func init() {
	slayCache = make(map[string][]*github.PullRequest)
}

type cmdSlayer struct {
	Token string
}

func (cmd *cmdSlayer) run(c *kingpin.ParseContext) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create in-cluster config")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create clientset")
	}

	pods, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list pods")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cmd.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientGithub := github.NewClient(tc)
	slayEmAll, err := getPodsToSlay(ctx, pods, clientGithub)
	if err != nil {
		return errors.Wrap(err, "failed to get pods to slay")
	}

	for _, pod := range slayEmAll {
		fmt.Printf("Slaying pod %s in namespace %s\n", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)
		err = clientset.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to slay pod")
		}
	}

	return nil
}

// getPodsToSlay gathers any pods with closed PRs.
func getPodsToSlay(ctx context.Context, pods *v1.PodList, clientGithub *github.Client) ([]v1.Pod, error) {
	var podsToSlay []v1.Pod

	for _, pod := range pods.Items {
		if _, ok := pod.Annotations[metadata.AnnotationCircleCIRepositoryName]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, metadata.AnnotationCircleCIRepositoryName)
			continue
		}
		if _, ok := pod.Annotations[metadata.AnnotationCircleCIRepositoryUsername]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, metadata.AnnotationCircleCIRepositoryUsername)
			continue
		}
		if _, ok := pod.Annotations[metadata.AnnotationCircleCIBranch]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, metadata.AnnotationCircleCIBranch)
			continue
		}

		prs, err := getClosedPrs(ctx, clientGithub, pod.Annotations[metadata.AnnotationCircleCIRepositoryUsername], pod.Annotations[metadata.AnnotationCircleCIRepositoryName])
		if err != nil {
			return nil, err
		}

		for _, pr := range prs {
			if pod.Annotations[metadata.AnnotationCircleCIBranch] == *pr.Head.Ref {
				podsToSlay = append(podsToSlay, pod)
			}
		}
	}

	return podsToSlay, nil
}

// getClosedPrs gathers all closed PRs for an owner and repo.
func getClosedPrs(ctx context.Context, client *github.Client, owner, repo string) ([]*github.PullRequest, error) {
	key := fmt.Sprintf("%s.%s", owner, repo)
	if val, ok := slayCache[key]; ok {
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

	slayCache[key] = allPrs
	return allPrs, nil
}

// Slayer declares the "slayer" sub command.
func Slayer(app *kingpin.Application) {
	c := new(cmdSlayer)

	cmd := app.Command("slayer", "Slay environments against closed PRs.").Action(c.run)
	cmd.Flag("token", "The Github access token.").Envar("GITHUB_TOKEN").Required().String()
}
