package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateSecret is a wrapper which will attempt to create and/or up a secret.
func SecretCreate(client *client.Clientset, new *v1.Secret) (*v1.Secret, error) {
	secret, err := client.CoreV1().Secrets(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.CoreV1().Secrets(new.ObjectMeta.Namespace).Update(new)
	}

	return secret, err
}
