package k8s

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	client *kubernetes.Clientset
	config *rest.Config
}