package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// DeploymentCreate is a wrapper which will attempt to create and/or up a deployment.
func DeploymentCreate(client *client.Clientset, new *extensions.Deployment) (*extensions.Deployment, error) {
	dply, err := client.Extensions().Deployments(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.Extensions().Deployments(new.ObjectMeta.Namespace).Update(new)
	}

	return dply, err
}
