package server

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server/k8s/env"
	"github.com/previousnext/m8s/server/k8s/utils"
	"k8s.io/client-go/kubernetes"
)

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

	err := stepClaims(srv.client, stream, srv.Namespace, srv.FilesystemSize)
	if err != nil {
		return err
	}

	err = stepService(srv.client, in, stream, srv.Namespace)
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

	return stepPod(srv.client, in, stream, srv.Namespace, srv.ApacheExporter)
}

// A step for provisioning caching storage.
func stepClaims(client *kubernetes.Clientset, stream pb.M8S_CreateServer, namespace, fs string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s PersistentVolumeClaim: Composer",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(namespace, env.CacheComposer, fs))
	if err != nil {
		return fmt.Errorf("failed to provision composer cache: %s", err)
	}

	err = stream.Send(&pb.CreateResponse{
		Message: "Creating K8s PersistentVolumeClaim: Yarn",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(namespace, env.CacheYarn, fs))
	if err != nil {
		return fmt.Errorf("failed to provision yarn cache: %s", err)
	}

	return nil
}

// A step to provision a Kubernetes service.
func stepService(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Service",
	})
	if err != nil {
		return err
	}

	_, err = utils.ServiceCreate(client, env.Service(namespace, in.Metadata.Name))
	if err != nil {
		return fmt.Errorf("failed create service: %s", err)
	}

	return nil
}

// A step to create a secret which contains http auth details.
func stepSecretBasicAuth(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace, name string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Secret: Basic Authentication",
	})
	if err != nil {
		return err
	}

	secret, err := env.Secret(namespace, name, in.Metadata.BasicAuth.User, in.Metadata.BasicAuth.Pass)
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
func stepIngress(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace, secret string) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Ingress",
	})
	if err != nil {
		return err
	}

	ing, err := env.Ingress(namespace, in.Metadata.Name, secret, in.Metadata.Domains)
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
func stepPod(client *kubernetes.Clientset, in *pb.CreateRequest, stream pb.M8S_CreateServer, namespace string, prom int32) error {
	err := stream.Send(&pb.CreateResponse{
		Message: "Creating K8s Pod",
	})
	if err != nil {
		return err
	}

	pod, err := env.Pod(namespace, in.Metadata.Name, in.GitCheckout.Repository, in.GitCheckout.Revision, in.Compose.Services, prom)
	if err != nil {
		return fmt.Errorf("failed build pod: %s", err)
	}

	_, err = utils.PodCreate(client, pod)
	if err != nil {
		return fmt.Errorf("failed create pod: %s", err)
	}

	return err
}
