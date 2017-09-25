package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// IngressCreate is a wrapper which will attempt to create and/or up an ingress.
func IngressCreate(client *client.Clientset, new *extensions.Ingress) (*extensions.Ingress, error) {
	ing, err := client.Extensions().Ingresses(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.Extensions().Ingresses(new.ObjectMeta.Namespace).Update(new)
	}

	return ing, err
}
