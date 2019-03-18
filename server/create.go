package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/previousnext/m8s/internal/podutils"
	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server/k8s/env"
	"github.com/previousnext/m8s/server/k8s/utils"
	"k8s.io/client-go/kubernetes"
)

// Create is used for creating a new environment.
func (srv Server) Create(in *pb.CreateRequest, stream pb.M8S_CreateServer) error {
	var authSecret string

	if in.Credentials.Token != srv.Token {
		return fmt.Errorf("token is incorrect")
	}

	if len(in.Compose.Services) < 1 {
		return fmt.Errorf("Docker Compose services not found")
	}

	if in.GitCheckout == nil {
		return fmt.Errorf("git checkout was not provided")
	}

	if in.GitCheckout.Revision == "" {
		return fmt.Errorf("git revision was not provided")
	}

	if in.GitCheckout.Repository == "" {
		return fmt.Errorf("git repository was not provided")
	}

	for _, cache := range srv.Cache.Directories {
		err := stepClaims(srv.client, stream, srv.Namespace, cache.Name, srv.Cache.Type, srv.Cache.Size)
		if err != nil {
			return err
		}
	}

	err := stepService(srv.client, in, stream, srv.Namespace)
	if err != nil {
		return err
	}

	if authProvided(in) {
		// Set the auth secret so the ingress can turn on auth at its layer.
		authSecret = fmt.Sprintf("%s-auth", in.Metadata.Name)

		err = stepSecretBasicAuth(srv.client, in, stream, srv.Namespace, authSecret)
		if err != nil {
			return err
		}
	}

	err = stepIngress(srv.client, in, stream, srv.Namespace, authSecret)
	if err != nil {
		return err
	}

	err = stepPod(srv.client, in, stream, srv.DockerCfg, srv.Namespace, srv.Cache.Directories)
	if err != nil {
		return err
	}

	return podutils.Tail(context.Background(), stream, srv.client, srv.Namespace, in.Metadata.Name)
}

// A step for provisioning caching storage.
func stepClaims(client *kubernetes.Clientset, stream pb.M8S_CreateServer, namespace, name, cacheType, cacheSize string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: fmt.Sprintf("Creating K8s PersistentVolumeClaim: %s\n", name),
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(env.PersistentVolumeClaimInput{
		Namespace: namespace,
		Name:      name,
		Type:      cacheType,
		Size:      cacheSize,
	}))
	if err != nil {
		return errors.Wrap(err, "failed to provision cache")
	}

	return nil
}

// A step to provision a Kubernetes service.
func stepService(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Service\n",
	})
	if err != nil {
		return err
	}

	svc, err := env.Service(env.ServiceInput{
		Namespace:   namespace,
		Name:        in.Metadata.Name,
		Annotations: in.Metadata.Annotations,
		Retention:   in.Metadata.Retention,
	})
	if err != nil {
		return err
	}

	_, err = utils.ServiceCreate(client, svc)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	return nil
}

// A step to create a secret which contains http auth details.
func stepSecretBasicAuth(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace, name string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Secret: Basic Authentication\n",
	})
	if err != nil {
		return err
	}

	secret, err := env.Secret(env.SecretInput{
		Namespace:   namespace,
		Name:        name,
		Annotations: in.Metadata.Annotations,
		User:        in.Metadata.BasicAuth.User,
		Pass:        in.Metadata.BasicAuth.Pass,
		Retention:   in.Metadata.Retention,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build secret")
	}

	_, err = utils.SecretCreate(client, secret)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	return nil
}

// A step to create an ingress for incoming traffic.
func stepIngress(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace, secret string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Ingress\n",
	})
	if err != nil {
		return err
	}

	ing, err := env.Ingress(env.IngressInput{
		Namespace:   namespace,
		Name:        in.Metadata.Name,
		Annotations: in.Metadata.Annotations,
		Secret:      secret,
		Retention:   in.Metadata.Retention,
		Domains:     in.Metadata.Domains,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build ingress")
	}

	// Create Basic Auth if required.
	_, err = utils.IngressCreate(client, ing)
	if err != nil {
		return errors.Wrap(err, "failed to create ingress")
	}

	return err
}

// A step for creating a pod (our Docker Compose environment).
func stepPod(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, dockercfg, namespace string, caches []CacheDirectory) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Pod\n",
	})
	if err != nil {
		return err
	}

	var inputCaches []env.PodInputCache

	for _, cache := range caches {
		inputCaches = append(inputCaches, env.PodInputCache{
			Name: cache.Name,
			Path: cache.Path,
		})
	}

	pod, err := env.Pod(env.PodInput{
		Namespace:       namespace,
		Name:            in.Metadata.Name,
		Annotations:     in.Metadata.Annotations,
		Repository:      in.GitCheckout.Repository,
		Revision:        in.GitCheckout.Revision,
		Retention:       in.Metadata.Retention,
		Services:        in.Compose.Services,
		Caches:          inputCaches,
		ImagePullSecret: dockercfg,
		Init:            in.Steps,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build pod")
	}

	_, err = utils.PodCreate(client, pod)
	if err != nil {
		return errors.Wrap(err, "failed to create pod")
	}

	return err
}
