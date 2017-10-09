package server

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const keyDockerCfg = ".dockercfg"

type Server struct {
	client            *kubernetes.Clientset
	config            *rest.Config
	Token             string
	Namespace         string
	FilesystemSize    string
	ApacheExporter    int32
	DockerCfgRegistry string
	DockerCfgUsername string
	DockerCfgPassword string
	DockerCfgEmail    string
	DockerCfgAuth     string
}

// DockerConfig is used mashalling Docker Configuration.
type DockerConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}
