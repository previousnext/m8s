package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client for executing.
type Client struct {
	client *kubernetes.Clientset
	config *rest.Config
}
