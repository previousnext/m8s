package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/previousnext/m8s/cmd/install/base"
	"github.com/previousnext/m8s/cmd/install/m8s"
	"github.com/previousnext/m8s/cmd/install/traefik"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type cmdInstall struct {
	KubeConfig string
	Token      string
	Namespace  string
	Domain     string
	Email      string
}

func (cmd *cmdInstall) run(c *kingpin.ParseContext) error {
	config, err := clientcmd.BuildConfigFromFlags("", cmd.KubeConfig)
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

	fmt.Println("Deployed!")
	fmt.Printf("Status: kubectl -n %s get pods\n", cmd.Namespace)
	fmt.Printf("Entrypoints: kubectl -n %s get svc\n", cmd.Namespace)

	return nil
}

// Install declares the "Install" sub command.
func Install(app *kingpin.Application) {
	c := new(cmdInstall)

	cmd := app.Command("install", "Installs M8s server components").Action(c.run)
	cmd.Flag("kubeconfig", "Path to users KubeConfig file").Default("$HOME/.kube/config").StringVar(&c.KubeConfig)
	cmd.Flag("namespace", "Namespace to deploy the M8s components into").Default("m8s").StringVar(&c.Namespace)
	cmd.Flag("token", "Token which the API server will use for authentication").Required().StringVar(&c.Token)
	cmd.Flag("domain", "Domain which the M8s api will respond on").Required().StringVar(&c.Domain)
	cmd.Flag("email", "Email to use for Let's Encrypt verification").Required().StringVar(&c.Email)
}
