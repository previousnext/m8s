package server

import (
	"encoding/json"
	"strings"

	"github.com/previousnext/m8s/server/k8s/env"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// Input is used for passing configuration for the cli to the server component.
type Input struct {
	Client    *kubernetes.Clientset
	Config    *rest.Config
	Token     string
	Namespace string
	SSH       string
	Cache     InputCache
	Dockercfg DockerRegistry
}

// InputCache is used as part of the input for the server.
type InputCache struct {
	Directories string
	Type        string
	Size        string
}

// New is used for returning a new M8s server.
func New(input Input) (Server, error) {
	srv := Server{
		client:    input.Client,
		config:    input.Config,
		Token:     input.Token,
		Namespace: input.Namespace,
		Docker:    input.Dockercfg,
		Cache: Cache{
			Type: input.Cache.Type,
			Size: input.Cache.Size,
		},
	}

	// Convert our cache directories into proper objects.
	cacheDirs, err := cacheDirectories(input.Cache.Directories)
	if err != nil {
		return srv, err
	}

	srv.Cache.Directories = cacheDirs

	err = dockercfgSync(input.Client, input.Namespace, input.Dockercfg)
	if err != nil {
		return srv, err
	}

	return srv, nil
}

// Helper function to sync Docker credentials.
func dockercfgSync(client *kubernetes.Clientset, namespace string, dockercfg DockerRegistry) error {
	auths := map[string]DockerCfg{
		dockercfg.Registry: {
			Username: dockercfg.Username,
			Password: dockercfg.Password,
			Email:    dockercfg.Email,
			Auth:     dockercfg.Auth,
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

// A helper function for converting a cache string into an object.
func cacheDirectories(arg string) ([]CacheDirectory, error) {
	var dirs []CacheDirectory

	for _, cache := range strings.Split(arg, ",") {
		sl := strings.Split(cache, ":")

		if len(sl) != 2 {
			continue
		}

		dirs = append(dirs, CacheDirectory{
			Name: sl[0],
			Path: sl[1],
		})
	}

	return dirs, nil
}
