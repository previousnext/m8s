package service

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"github.com/pkg/errors"
)

// Deploy will create the Service if not present.
func Deploy(client *kubernetes.Clientset, service *corev1.Service) error {
	_, err := client.CoreV1().Services(service.ObjectMeta.Namespace).Create(service)
	// We don't do anything if this is an existing resource.
	if kerrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
