package base

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// Install will setup the base components for M8s.
func Install(client *kubernetes.Clientset, namespace string) error {
	fmt.Println("Installing: Namespace")

	return installNamespace(client, namespace)
}
