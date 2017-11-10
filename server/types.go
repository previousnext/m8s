package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const keyDockerCfg = ".dockercfg"

// Server is the M8s server for running builds.
type Server struct {
	client     *kubernetes.Clientset
	config     *rest.Config
	Token      string
	Namespace  string
	SSHService string
	Cache      Cache
	Docker     DockerRegistry
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

// DockerRegistry contains Docker Hub credentials and registry information.
type DockerRegistry struct {
	Registry string
	Username string
	Password string
	Email    string
	Auth     string
}

// DockerCfg is used mashalling Docker Configuration.
type DockerCfg struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}
