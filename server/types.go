package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const keyDockerCfg = ".dockercfg"

// Server is the M8s server for running builds.
type Server struct {
	client         *kubernetes.Clientset
	config         *rest.Config
	Token          string
	Namespace      string
	FilesystemSize string
	SSHService     string
	ApacheExporter int32
	Docker         DockerRegistry
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
