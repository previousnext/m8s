package namespace

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// Deploy will create the Namespace if not present.
func Deploy(client *kubernetes.Clientset, namespace *corev1.Namespace) error {
	_, err := client.CoreV1().Namespaces().Create(namespace)
	// We don't do anything if this is an existing resource.
	if kerrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
