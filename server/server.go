package server

import (
	"encoding/json"

	"github.com/previousnext/m8s/server/k8s/env"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// New is used for returning a new M8s server.
func New(client *kubernetes.Clientset, config *rest.Config, token, namespace, fs string, exporter int32, dockerreg, dockeruser, dockerpass, dockeremail, dockerauth string) (Server, error) {
	srv := Server{
		client:            client,
		config:            config,
		Token:             token,
		Namespace:         namespace,
		FilesystemSize:    fs,
		ApacheExporter:    exporter,
		DockerCfgRegistry: dockerreg,
		DockerCfgUsername: dockeruser,
		DockerCfgPassword: dockerpass,
		DockerCfgEmail:    dockeremail,
		DockerCfgAuth:     dockerauth,
	}

	err := dockercfgSync(client, namespace, dockerreg, dockeruser, dockerpass, dockeremail, dockerauth)
	if err != nil {
		return srv, err
	}

	return srv, nil
}

// Helper function to sync Docker credentials.
func dockercfgSync(client *kubernetes.Clientset, namespace, registry, username, password, email, auth string) error {
	auths := map[string]DockerConfig{
		registry: {
			Username: username,
			Password: password,
			Email:    email,
			Auth:     auth,
		},
	}

	dockerconfig, err := json.Marshal(auths)
	if err != nil {
		return err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      env.SecretDockerCfg,
		},
		Data: map[string][]byte{
			keyDockerCfg: dockerconfig,
		},
		Type: v1.SecretTypeDockercfg,
	}

	_, err = utils.SecretCreate(client, secret)
	if err != nil {
		return err
	}

	return nil
}
