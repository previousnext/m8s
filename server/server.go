package server

import (
	"strings"

	"k8s.io/client-go/kubernetes"
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

	return srv, nil
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
