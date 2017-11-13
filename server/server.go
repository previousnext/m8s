package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// New is used for returning a new M8s server.
func New(client *kubernetes.Clientset, config *rest.Config, token, namespace, cacheType, cacheSize string, exporter int32) (Server, error) {
	srv := Server{
		client:         client,
		config:         config,
		Token:          token,
		Namespace:      namespace,
		CacheType:      cacheType,
		CacheSize:      cacheSize,
		ApacheExporter: exporter,
	}

	return srv, nil
}
