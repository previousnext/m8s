package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// IngressCreate is a wrapper which will attempt to create and/or up an ingress.
func IngressCreate(client *kubernetes.Clientset, new *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	ing, err := client.Ingresses(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.Ingresses(new.ObjectMeta.Namespace).Update(new)
	}

	return ing, err
}
