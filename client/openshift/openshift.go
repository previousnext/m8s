package openshift

import (
	"github.com/mitchellh/go-homedir"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/client/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const Name = "openshift"

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

	client.client = c
	client.config = cfg

	r, err := routev1.NewForConfig(cfg)
	if err != nil {
		return client, errors.Wrap(err, "failed to get Openshift client")
	}

	client.routev1client = r

	return client, nil
}
