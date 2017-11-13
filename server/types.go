package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Server is the M8s server for running builds.
type Server struct {
	client     *kubernetes.Clientset
	config     *rest.Config
	Token      string
	Namespace  string
	SSHService string
	Cache      Cache
}

// Cache decares all the cached paths with will be a path of the builds.
type Cache struct {
	Directories []CacheDirectory
	Type        string
	Size        string
}

// CacheDirectory is used for caching a path on a pod.
type CacheDirectory struct {
	Name string
	Path string
}
