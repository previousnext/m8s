package configmap

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// Deploy will create the ConfigMap if not present.
func Deploy(client *kubernetes.Clientset, configmap *corev1.ConfigMap) error {
	_, err := client.CoreV1().ConfigMaps(configmap.ObjectMeta.Namespace).Create(configmap)
	// We don't do anything if this is an existing resource.
	if kerrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
