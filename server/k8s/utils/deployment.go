package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
)

// DeploymentCreate is a wrapper which will attempt to create and/or up a deployment.
func DeploymentCreate(client *kubernetes.Clientset, new *appsv1.Deployment) (*appsv1.Deployment, error) {
	dply, err := client.AppsV1().Deployments(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.AppsV1().Deployments(new.ObjectMeta.Namespace).Update(new)
	}

	return dply, err
}
