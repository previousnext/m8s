package main

import (
	"k8s.io/client-go/rest"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

const keyDockerCfg = ".dockercfg"

type server struct {
	client *client.Clientset
	config *rest.Config
}

// DockerConfig is used mashalling Docker Configuration.
type DockerConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}
