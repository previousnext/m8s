package m8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// Install will setup all the M8s components.
func Install(client *kubernetes.Clientset, namespace, token, letsEncryptDomain, letsEncryptEmail string) error {
	fmt.Println("Installing M8s API: Deployment")

	err := installDeployment(client, namespace, token, letsEncryptDomain, letsEncryptEmail)
	if err != nil {
		return err
	}

	fmt.Println("Installing M8s API: PVC")

	err = installPVC(client, namespace)
	if err != nil {
		return err
	}

	fmt.Println("Installing M8s API: Secret")

	err = installSecret(client, namespace)
	if err != nil {
		return err
	}

	fmt.Println("Installing M8s API: Servce")

	return installService(client, namespace)
}
