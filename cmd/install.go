package cmd

import (
	"html/template"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/cmd/install/base"
	"github.com/previousnext/m8s/cmd/install/m8s"
	"github.com/previousnext/m8s/cmd/install/traefik"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const success = `Deployed!

To get the status of the components:

  kubectl -n {{ . }} get pods

To get a list of IP address for setting up DNS:

  kubectl -n {{ . }} get svc`

type cmdInstall struct {
	APIServer  string
	KubeConfig string
	Token      string
	Namespace  string
	Domain     string
	Email      string
}

func (cmd *cmdInstall) run(c *kingpin.ParseContext) error {
	kubeconfig, err := homedir.Expand(cmd.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "failed to get kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags(cmd.APIServer, kubeconfig)
	if err != nil {
		return errors.Wrap(err, "failed to get clientcmd config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes client")
	}

	err = base.Install(client, cmd.Namespace)
	if err != nil {
		return errors.Wrap(err, "failed to install Base")
	}

	err = traefik.Install(client, cmd.Namespace)
	if err != nil {
		return errors.Wrap(err, "failed to install Traefik")
	}

	err = m8s.Install(client, cmd.Namespace, cmd.Token, cmd.Domain, cmd.Email)
	if err != nil {
		return errors.Wrap(err, "failed to install M8s API")
	}

	t := template.Must(template.New("success").Parse(success))

	err = t.Execute(os.Stdout, cmd.Namespace)
	if err != nil {
		return errors.Wrap(err, "failed to generate success message")
	}

	return nil
}

// Install declares the "Install" sub command.
func Install(app *kingpin.Application) {
	c := new(cmdInstall)

	cmd := app.Command("install", "Installs M8s server components").Action(c.run)
	cmd.Flag("apiserver", "Kubernetes apiserver endpoint eg. http://localhost:8080").StringVar(&c.APIServer)
	cmd.Flag("kubeconfig", "Path to users KubeConfig file").Default("~/.kube/config").StringVar(&c.KubeConfig)
	cmd.Flag("namespace", "Namespace to deploy the M8s components into").Default("m8s").StringVar(&c.Namespace)
	cmd.Flag("token", "Token which the API server will use for authentication").Required().StringVar(&c.Token)
	cmd.Flag("domain", "Domain which the M8s api will respond on").Required().StringVar(&c.Domain)
	cmd.Flag("email", "Email to use for Let's Encrypt verification").Required().StringVar(&c.Email)
}
