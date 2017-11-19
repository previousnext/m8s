package traefik

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// Install will setup the Traefik routing Ingress layer.
func Install(client *kubernetes.Clientset, namespace string) error {
	fmt.Println("Installing Traefik: Deployment")

	err := installDeployment(client, namespace)
	if err != nil {
		return err
	}

	fmt.Println("Installing Traefik: Service")

	return installService(client, namespace)
}
