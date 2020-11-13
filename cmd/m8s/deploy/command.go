package deploy

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/previousnext/m8s/internal/environ"
	"github.com/previousnext/m8s/internal/k8s/pod/annotations"
	podbuilder "github.com/previousnext/m8s/internal/k8s/pod/builder"
)

type command struct {
	// Name which will be associated with this environment.
	Name string
	// Namespace which this Pod will be created in
	Namespace string
	// Domain which will be associated with this environment.
	Domain string
	// GitRepository which will be cloned during an environment build.
	GitRepository string
	// GitRevision which will be cloned during an environment build.
	GitRevision string
	// Config file to be loaded.
	Configs []string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	// Load the Docker Compose file, we are going to use alot of its
	// configuration for this build.
	//dc, err := config.Load(cmd.DockerCompose)
	//if err != nil {
	//	return fmt.Errorf("failed to load Docker Compose configuration: %w", err)
	//}

	pod, err := podbuilder.Generate(podbuilder.GenerateParams{
		Name: cmd.Name,
		Namespace: cmd.Namespace,
		Domain: cmd.Domain,
		ExtraAnnotations: annotations.FromEnvironment(os.Environ()),
		// These are additional environment variables that have been provided outside of this build, with the intent
		// for them to be injected into our running containers.
		//   eg. M8S_ENV_FOO=bar, will inject FOO=bar into the containers.
		ExtraEnvironment: environ.GetWithPrefix("M8S_ENV_", os.Environ()),
	})
	if err != nil {
		return fmt.Errorf("unable to generate Pod: %w", err)
	}

	raw, err := json.Marshal(pod)
	if err != nil {
		return fmt.Errorf("unable to generate Pod: %w", err)
	}

	fmt.Println(string(raw))

	/*pod, err = podclient.Create(nil, pod)
	if err != nil {
		return fmt.Errorf("failed to create Pod: %w", err)
	}*/

	return nil

	// return podlogs.Stream(context.Background())
}

// Command declares the "deploy" sub command.
func Command(app *kingpin.Application) {
	c := new(command)

	cmd := app.Command("deploy", "Deploy the temporary environment").Action(c.run)

	cmd.Flag("git-repository", "Git repository to clone from").Required().Envar("M8S_GIT_REPO").StringVar(&c.GitRepository)
	cmd.Flag("git-revision", "Git revision to checkout during clone").Required().Envar("M8S_GIT_REVISION").StringVar(&c.GitRevision)

	cmd.Flag("namespace", "Namespace which this Pod will be created in").Envar("M8S_NAMESPACE").Default(corev1.NamespaceDefault).StringVar(&c.Namespace)

	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Arg("domain", "Domain for this environment to run on").Required().StringVar(&c.Domain)
	cmd.Arg("config", "Config file for M8s").Required().StringsVar(&c.Configs)
}
