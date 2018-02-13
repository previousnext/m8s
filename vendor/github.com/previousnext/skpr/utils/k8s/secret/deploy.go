package secret

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"github.com/pkg/errors"
)

// Deploy will create the Secret and fallback to updating if it exists.
func Deploy(client *kubernetes.Clientset, secret *corev1.Secret) error {
	_, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(secret)
	if kerrors.IsAlreadyExists(err) {
		_, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Update(secret)
		if err != nil {
			return errors.Wrap(err, "failed to update")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
