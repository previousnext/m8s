package cmd

import (
	"context"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/cmd/metadata"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	slayCache map[string][]*github.PullRequest
)

func init() {
	slayCache = make(map[string][]*github.PullRequest)
}

type cmdSlayer struct {
	MasterURL string
	Kubeconfig string

	Token     string
	Namespace string

	Grace int64

	Debug bool
}

func (cmd *cmdSlayer) run(c *kingpin.ParseContext) error {
	ctx := context.Background()

	config, err := clientcmd.BuildConfigFromFlags(cmd.MasterURL, cmd.Kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	pods, err := clientset.CoreV1().Pods(cmd.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list pods")
	}

	clientGithub := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cmd.Token},
	)))

	slayEmAll, err := getPodsToSlay(ctx, pods, clientGithub)
	if err != nil {
		return fmt.Errorf("failed to get Pods: %w", err)
	}

	for _, pod := range slayEmAll {
		fmt.Printf("Slaying pod: %s/%s\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)

		if cmd.Debug {
			continue
		}

		err = clientset.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{
			GracePeriodSeconds: &cmd.Grace,
		})
		if err != nil {
			return fmt.Errorf("failed to slay Pods (it got away!): %w", err)
		}
	}

	return nil
}

// getPodsToSlay gathers any pods with closed PRs.
func getPodsToSlay(ctx context.Context, pods *v1.PodList, clientGithub *github.Client) ([]v1.Pod, error) {
	var podsToSlay []v1.Pod

	createdGracePeriod := time.Now().Add(time.Duration(-12) * time.Hour)

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
		if _, ok := pod.Annotations[metadata.AnnotationCircleCISHA1]; !ok {
			fmt.Printf("Pod %s missing the %s annotation, skipping\n", pod.ObjectMeta.Name, metadata.AnnotationCircleCISHA1)
			continue
		}

		prs, err := getOpenPRs(ctx, clientGithub, pod.Annotations[metadata.AnnotationCircleCIRepositoryUsername], pod.Annotations[metadata.AnnotationCircleCIRepositoryName])
		if err != nil {
			return nil, err
		}

		if isOpenPR(prs, pod.Annotations[metadata.AnnotationCircleCIBranch], pod.Annotations[metadata.AnnotationCircleCISHA1]) {
			fmt.Printf("Pod %s has an open pull request, skipping\n", pod.ObjectMeta.Name)
			continue
		}

		fmt.Printf("Comparing Pod %s/%s with created time %s to grace period %s\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, pod.ObjectMeta.CreationTimestamp.String(), createdGracePeriod.String())

		if pod.ObjectMeta.CreationTimestamp.After(createdGracePeriod) {
			fmt.Printf("Pod %s was only just created and might not have a pull request yet, skipping\n", pod.ObjectMeta.Name)
			continue
		}

		podsToSlay = append(podsToSlay, pod)
	}

	return podsToSlay, nil
}

// Helper function to check if a branch is on the list.
func isOpenPR(list []*github.PullRequest, branch, sha string) bool {
	for _, item := range list {
		if branch == *item.Head.Ref && sha == *item.Head.SHA {
			return true
		}
	}

	return false
}

// getOpenPRs gathers all open PRs for an owner and repo.
func getOpenPRs(ctx context.Context, client *github.Client, owner, repo string) ([]*github.PullRequest, error) {
	key := fmt.Sprintf("%s.%s", owner, repo)
	if val, ok := slayCache[key]; ok {
		fmt.Printf("Loading PRs from cache for owner %s and repo %s\n", owner, repo)
		return val, nil
	}

	var allPrs []*github.PullRequest
	opt := &github.PullRequestListOptions{
		State:       "open",
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

	cmd.Flag("master", "Location of the Kubernetes master.").Envar("SLAYER_MASTER").StringVar(&c.MasterURL)
	cmd.Flag("kubeconfig", "Location of the Kubernetes Kubeconfig.").Envar("SLAYER_KUBECONFIG").StringVar(&c.Kubeconfig)

	cmd.Flag("token", "The Github access token.").Envar("GITHUB_TOKEN").Required().StringVar(&c.Token)
	cmd.Flag("namespace", "The Kubernetes namespace to slay pods in.").Default(v1.NamespaceAll).Envar("SLAYER_NAMESPACE").StringVar(&c.Namespace)

	cmd.Flag("grace", "How long a Pod should be allowed to shutdown").Envar("SLAYER_GRACE").Int64Var(&c.Grace)

	cmd.Flag("debug", "Pods will not be terminated").Envar("SLAYER_DEBUG").BoolVar(&c.Debug)
}
