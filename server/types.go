package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Server is the M8s server for running builds.
type Server struct {
	client         *kubernetes.Clientset
	config         *rest.Config
	Token          string
	Namespace      string
	CacheType      string
	CacheSize      string
	SSHService     string
	ApacheExporter int32
}
