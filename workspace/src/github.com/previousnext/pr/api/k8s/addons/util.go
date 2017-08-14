package addons

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// Helper function for creating deployments.
func createDeployment(client *client.Clientset, dply *v1beta1.Deployment) error {
	_, err := client.Extensions().Deployments(dply.ObjectMeta.Namespace).Create(dply)
	if errors.IsAlreadyExists(err) {
		_, err = client.Extensions().Deployments(dply.ObjectMeta.Namespace).Update(dply)
		if err != nil {
			return err
		}
	}

	return nil
}
