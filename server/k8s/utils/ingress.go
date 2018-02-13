package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

// IngressCreate is a wrapper which will attempt to create and/or up an ingress.
func IngressCreate(client *kubernetes.Clientset, new *extensionsv1beta1.Ingress) (*extensionsv1beta1.Ingress, error) {
	ing, err := client.ExtensionsV1beta1().Ingresses(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.ExtensionsV1beta1().Ingresses(new.ObjectMeta.Namespace).Update(new)
	}

	return ing, err
}
