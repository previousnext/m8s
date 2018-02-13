package pvc

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"github.com/pkg/errors"
)

// Deploy will create the PersistentVolumeClaim if not present.
func Deploy(client *kubernetes.Clientset, pvc *corev1.PersistentVolumeClaim) error {
	_, err := client.CoreV1().PersistentVolumeClaims(pvc.ObjectMeta.Namespace).Create(pvc)
	// We don't do anything if this is an existing resource.
	if kerrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
