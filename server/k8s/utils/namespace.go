package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

// NamespaceCreate is a wrapper which will only create a namespace.
func NamespaceCreate(client *kubernetes.Clientset, new *corev1.Namespace) (*corev1.Namespace, error) {
	ns, err := client.CoreV1().Namespaces().Create(new)
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, err
	}

	return ns, nil
}
