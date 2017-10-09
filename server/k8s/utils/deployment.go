package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

// DeploymentCreate is a wrapper which will attempt to create and/or up a deployment.
func DeploymentCreate(client *kubernetes.Clientset, new *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	dply, err := client.Apps().Deployments(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.Apps().Deployments(new.ObjectMeta.Namespace).Update(new)
	}

	return dply, err
}
