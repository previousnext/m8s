package k8s

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/client/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Name of the client.
const Name = "k8s"

// New client implemented on top of Kubernetes.
func New(params types.ClientParams) (types.Client, error) {
	var client Client

	kubeconfig, err := homedir.Expand(params.KubeConfig)
	if err != nil {
		return client, errors.Wrap(err, "failed to get kubeconfig homedir")
	}

	cfg, err := clientcmd.BuildConfigFromFlags(params.Master, kubeconfig)
	if err != nil {
		return client, errors.Wrap(err, "failed to get Kubernetes config")
	}

	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return client, errors.Wrap(err, "failed to get Kubernetes client")
	}

	return Client{client: c, config: cfg}, nil
}
