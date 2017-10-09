package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// ServiceCreate is a wrapper which will attempt to create and/or up a service.
func ServiceCreate(client *kubernetes.Clientset, new *v1.Service) (*v1.Service, error) {
	svc, err := client.CoreV1().Services(new.ObjectMeta.Namespace).Create(new)
	// We don't do anything if this is an existing resource.
	if errors.IsAlreadyExists(err) {
		return svc, nil
	}

	return svc, err
}
