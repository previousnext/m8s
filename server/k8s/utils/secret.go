package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

// SecretCreate is a wrapper which will attempt to create and/or up a secret.
func SecretCreate(client *kubernetes.Clientset, new *corev1.Secret) (*corev1.Secret, error) {
	secret, err := client.CoreV1().Secrets(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.CoreV1().Secrets(new.ObjectMeta.Namespace).Update(new)
	}

	return secret, err
}
