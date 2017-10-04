package main

import (
	"fmt"

	"github.com/previousnext/m8s/api/k8s/env"
	"github.com/previousnext/m8s/api/k8s/utils"
	pb "github.com/previousnext/m8s/pb"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

func (srv server) Create(in *pb.CreateRequest, stream pb.M8S_CreateServer) error {
	var authSecret string

	if in.Credentials.Token != *cliToken {
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

	err := stepClaims(srv.client, stream)
	if err != nil {
		return err
	}

	err = stepService(srv.client, in, stream)
	if err != nil {
		return err
	}

	if authProvided(in) {
		// Set the auth secret so the ingress can turn on auth at its layer.
		authSecret = fmt.Sprintf("%s-auth", in.Metadata.Name)

		err = stepSecretBasicAuth(srv.client, in, stream, authSecret)
		if err != nil {
			return err
		}
	}

	err = stepIngress(srv.client, in, stream, authSecret)
	if err != nil {
		return err
	}

	return stepPod(srv.client, in, stream)
}

// A step for provisioning caching storage.
func stepClaims(client *client.Clientset, stream pb.M8S_CreateServer) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s PersistentVolumeClaim: Composer",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(*cliNamespace, env.CacheComposer, *cliFilesystemSize))
	if err != nil {
		return fmt.Errorf("failed to provision composer cache: %s", err)
	}

	err = stream.Send(&pb.CreateResponse{
		Message: "Creating K8s PersistentVolumeClaim: Yarn",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(*cliNamespace, env.CacheYarn, *cliFilesystemSize))
	if err != nil {
		return fmt.Errorf("failed to provision yarn cache: %s", err)
	}

	return nil
}

// A step to provision a Kubernetes service.
func stepService(client *client.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Service",
	})
	if err != nil {
		return err
	}

	_, err = utils.ServiceCreate(client, env.Service(*cliNamespace, in.Metadata.Name))
	if err != nil {
		return fmt.Errorf("failed create service: %s", err)
	}

	return nil
}

// A step to create a secret which contains http auth details.
func stepSecretBasicAuth(client *client.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, name string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Secret: Basic Authentication",
	})
	if err != nil {
		return err
	}

	secret, err := env.Secret(*cliNamespace, name, in.Metadata.BasicAuth.User, in.Metadata.BasicAuth.Pass)
	if err != nil {
		return fmt.Errorf("failed build secret: %s", err)
	}

	_, err = utils.SecretCreate(client, secret)
	if err != nil {
		return fmt.Errorf("failed create service: %s", err)
	}

	return nil
}

// A step to create an ingress for incoming traffic.
func stepIngress(client *client.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, secret string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Ingress",
	})
	if err != nil {
		return err
	}

	ing, err := env.Ingress(*cliNamespace, in.Metadata.Name, secret, in.Metadata.Domains)
	if err != nil {
		return fmt.Errorf("failed build ingress: %s", err)
	}

	// Create Basic Auth if required.
	_, err = utils.IngressCreate(client, ing)
	if err != nil {
		return fmt.Errorf("failed create ingress: %s", err)
	}

	return err
}

// A step for creating a pod (our Docker Compose environment).
func stepPod(client *client.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Pod",
	})
	if err != nil {
		return err
	}

	pod, err := env.Pod(*cliNamespace, in.Metadata.Name, in.GitCheckout.Repository, in.GitCheckout.Revision, in.Compose.Services, *cliPrometheusApache)
	if err != nil {
		return fmt.Errorf("failed build pod: %s", err)
	}

	_, err = utils.PodCreate(client, pod)
	if err != nil {
		return fmt.Errorf("failed create pod: %s", err)
	}

	return err
}
