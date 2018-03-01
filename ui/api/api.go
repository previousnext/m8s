package api

import (
	k8sapi "github.com/previousnext/m8s/ui/api/k8s"
	mockapi "github.com/previousnext/m8s/ui/api/mock"
)

// New returns a new API server.
func New(master, config, namespace string, mock bool) API {
	if mock {
		return mockapi.Server{}
	}

	return k8sapi.Server{
		Namespace: namespace,
		Master:    master,
		Config:    config,
	}
}
