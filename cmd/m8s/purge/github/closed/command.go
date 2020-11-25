package closed

import (
	"fmt"
	"context"

	corev1 "k8s.io/api/core/v1"
	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/previousnext/m8s/internal/k8s/pod/annotations"
)

var (
	// Cache of pull requests for a repository.
	cache map[string][]*github.PullRequest
)

func init() {
	cache = make(map[string][]*github.PullRequest)
}

type command struct {
	MasterURL string
	Kubeconfig string
	Namespace string
	Grace int64

	GithubToken string

	DryRun bool
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	ctx := context.Background()

	config, err := clientcmd.BuildConfigFromFlags(cmd.MasterURL, cmd.Kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	list, err := clientset.CoreV1().Pods(cmd.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list pods")
	}

	gh := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cmd.GithubToken},
	)))

	pods, err := getPodsToDelete(ctx, gh, list)
	if err != nil {
		return fmt.Errorf("failed to delete Pods: %w", err)
	}

	for _, pod := range pods {
		fmt.Printf("Deleting pod: %s/%s\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)

		if cmd.DryRun {
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

// Helper function to get Pods which will be deleted.
func getPodsToDelete(ctx context.Context, gh *github.Client, list *corev1.PodList) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	for _, pod := range list.Items {
		repoUser, repoName, err := getRepositoryFromAnnotations(pod)
		if err != nil {
			// @todo, Log it.
			continue
		}

		branch, sha, err := getBranchFromAnnotations(pod)
		if err != nil {
			// @todo, Log it.
			continue
		}

		prs, err := getPullRequests(ctx, gh, repoUser, repoName)
		if err != nil {
			return nil, err
		}

		if !exists(prs, branch, sha) {
			fmt.Println("Skipping Pod because it is not on the list:", pod.ObjectMeta.Name)
			continue
		}

		pods = append(pods, pod)
	}

	return pods, nil
}

func getRepositoryFromAnnotations(pod corev1.Pod) (string, string, error) {
	if _, ok := pod.Annotations[annotations.AnnotationCircleCIRepositoryUsername]; !ok {
		return "", "", fmt.Errorf("%s missing the %s annotation, skipping", pod.ObjectMeta.Name, annotations.AnnotationCircleCIRepositoryUsername)
	}

	if _, ok := pod.Annotations[annotations.AnnotationCircleCIRepositoryName]; !ok {
		return "", "", fmt.Errorf("%s missing the %s annotation, skipping", pod.ObjectMeta.Name, annotations.AnnotationCircleCIRepositoryName)
	}

	return pod.Annotations[annotations.AnnotationCircleCIRepositoryUsername], pod.Annotations[annotations.AnnotationCircleCIRepositoryName], nil
}

func getBranchFromAnnotations(pod corev1.Pod) (string, string, error) {
	if _, ok := pod.Annotations[annotations.AnnotationCircleCIBranch]; !ok {
		return "", "", fmt.Errorf("%s missing the %s annotation, skipping", pod.ObjectMeta.Name, annotations.AnnotationCircleCIBranch)
	}

	if _, ok := pod.Annotations[annotations.AnnotationCircleCISHA1]; !ok {
		return "", "", fmt.Errorf("%s missing the %s annotation, skipping", pod.ObjectMeta.Name, annotations.AnnotationCircleCISHA1)
	}

	return pod.Annotations[annotations.AnnotationCircleCIBranch], pod.Annotations[annotations.AnnotationCircleCISHA1], nil
}

// Helper function to get a list of pull requests with a given state.
func getPullRequests(ctx context.Context, client *github.Client, owner, repo string) ([]*github.PullRequest, error) {
	key := fmt.Sprintf("%s.%s", owner, repo)

	if val, ok := cache[key]; ok {
		return val, nil
	}

	var all []*github.PullRequest

	opt := &github.PullRequestListOptions{
		State:       "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		all = append(all, prs...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	cache[key] = all

	return all, nil
}

// Helper function to check if a branch is on the list.
func exists(list []*github.PullRequest, branch, sha string) bool {
	for _, item := range list {
		if branch == *item.Head.Ref && sha == *item.Head.SHA {
			return true
		}
	}

	return false
}

// Command declares the "closed" sub command.
func Command(app *kingpin.CmdClause) {
	c := new(command)

	cmd := app.Command("closed", "Purge environments which are Github Pull Requests which are closed").Action(c.run)

	cmd.Flag("master", "Location of the Kubernetes master.").Envar("M8S_PURGE_MASTER").StringVar(&c.MasterURL)
	cmd.Flag("kubeconfig", "Location of the Kubernetes Kubeconfig.").Envar("M8S_PURGE_KUBECONFIG").StringVar(&c.Kubeconfig)
	cmd.Flag("namespace", "The Kubernetes namespace to slay pods in.").Default(corev1.NamespaceAll).Envar("M8S_PURGE_NAMESPACE").StringVar(&c.Namespace)
	cmd.Flag("grace", "How long a Pod should be allowed to shutdown").Envar("M8S_PURGE_GRACE").Int64Var(&c.Grace)

	cmd.Flag("github-token", "The Github access token.").Envar("M8S_PURGE_GITHUB_TOKEN").Required().StringVar(&c.GithubToken)

	cmd.Flag("dry-run", "Print out what will happen").Envar("M8S_PURGE_DRY_RUN").BoolVar(&c.DryRun)
}
