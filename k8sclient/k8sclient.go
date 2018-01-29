package k8sclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// New returns a new Kubernets client and config.
func New(master, kubeconfig string) (*kubernetes.Clientset, *rest.Config, error) {
	config, err := getConfig(master, kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return client, config, nil
}

// Helper function for inspecting if the user is incluster or outside.
func getConfig(master, kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags(master, kubeconfig)
	}

	return rest.InClusterConfig()
}
