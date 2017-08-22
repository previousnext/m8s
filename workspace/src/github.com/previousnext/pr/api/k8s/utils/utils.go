package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateDeployment is a wrapper which will attempt to create and/or up a deployment.
func CreateDeployment(client *client.Clientset, dply *v1beta1.Deployment) error {
	_, err := client.Extensions().Deployments(dply.ObjectMeta.Namespace).Create(dply)
	if errors.IsAlreadyExists(err) {
		_, err = client.Extensions().Deployments(dply.ObjectMeta.Namespace).Update(dply)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateService is a wrapper which will attempt to create and/or up a service.
func CreateService(client *client.Clientset, svc *v1.Service) error {
	_, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)
	// We don't do anything if this is an existing resource.
	if errors.IsAlreadyExists(err) {
		return nil
		// We still need to tell the admin if this is an error on an existing object.
	} else if err != nil {
		return err
	}

	return nil
}
