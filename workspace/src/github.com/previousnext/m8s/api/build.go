package main

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/previousnext/m8s/api/k8s/env"
	"github.com/previousnext/m8s/api/k8s/utils"
	pb "github.com/previousnext/m8s/pb"
	"k8s.io/client-go/rest"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

func (srv server) Build(in *pb.BuildRequest, stream pb.M8S_BuildServer) error {
	if in.Credentials.Token != *cliToken {
		return fmt.Errorf("token is incorrect")
	}

	if in.Metadata.BasicAuth == nil {
		return fmt.Errorf("basic auth was not provided")
	}

	if in.Metadata.BasicAuth.User == "" {
		return fmt.Errorf("basic auth field was not provided: user")
	}

	if in.Metadata.BasicAuth.Pass == "" {
		return fmt.Errorf("basic auth field was not provided: pass")
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

	if in.Exec.Container == "" {
		return fmt.Errorf("exec container was not provided")
	}

	if len(in.Exec.Steps) < 1 {
		return fmt.Errorf("exec steps were not provided")
	}

	if in.Keep == "" {
		return fmt.Errorf("keep duration was not provided")
	}

	keep, err := time.ParseDuration(in.Keep)
	if err != nil {
		return fmt.Errorf("failed marshall field 'keep': %s", err)
	}

	var (
		blackDeath = time.Now().Unix() + keep.Nanoseconds()
		authSecret = fmt.Sprintf("%s-auth", in.Metadata.Name)
	)

	err = stepClaims(srv.client, stream)
	if err != nil {
		return err
	}

	err = stepService(srv.client, in, stream, blackDeath)
	if err != nil {
		return err
	}

	err = stepSecret(srv.client, in, stream, blackDeath, authSecret)
	if err != nil {
		return err
	}

	err = stepIngress(srv.client, in, stream, blackDeath, authSecret)
	if err != nil {
		return err
	}

	err = stepPod(srv.client, in, stream, blackDeath)
	if err != nil {
		return err
	}

	return stepCommands(srv.client, in, stream, srv.config, *cliNamespace)
}

// A step for provisioning caching storage.
func stepClaims(client *client.Clientset, stream pb.M8S_BuildServer) error {
	err := stream.Send(&pb.BuildResponse{
		Message: "Creating K8s PersistentVolumeClaim: Composer",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(*cliNamespace, env.CacheComposer, *cliCacheSize))
	if err != nil {
		return fmt.Errorf("failed to provision composer cache: %s", err)
	}

	err = stream.Send(&pb.BuildResponse{
		Message: "Creating K8s PersistentVolumeClaim: Yarn",
	})
	if err != nil {
		return err
	}

	_, err = utils.PersistentVolumeClaimCreate(client, env.PersistentVolumeClaim(*cliNamespace, env.CacheYarn, *cliCacheSize))
	if err != nil {
		return fmt.Errorf("failed to provision yarn cache: %s", err)
	}

	return nil
}

// A step to provision a Kubernetes service.
func stepService(client *client.Clientset, in *pb.BuildRequest, stream pb.M8S_BuildServer, blackDeath int64) error {
	err := stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Service",
	})
	if err != nil {
		return err
	}

	_, err = utils.ServiceCreate(client, env.Service(blackDeath, *cliNamespace, in.Metadata.Name))
	if err != nil {
		return fmt.Errorf("failed create service: %s", err)
	}

	return nil
}

// A step to create a secret which contains http auth details.
func stepSecret(client *client.Clientset, in *pb.BuildRequest, stream pb.M8S_BuildServer, blackDeath int64, name string) error {
	err := stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Secret: auth",
	})
	if err != nil {
		return err
	}

	secret, err := env.Secret(blackDeath, *cliNamespace, name, in.Metadata.BasicAuth.User, in.Metadata.BasicAuth.Pass)
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
func stepIngress(client *client.Clientset, in *pb.BuildRequest, stream pb.M8S_BuildServer, blackDeath int64, secret string) error {
	err := stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Ingress",
	})
	if err != nil {
		return err
	}

	ing, err := env.Ingress(blackDeath, *cliNamespace, in.Metadata.Name, secret, in.Metadata.Domains)
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
func stepPod(client *client.Clientset, in *pb.BuildRequest, stream pb.M8S_BuildServer, blackDeath int64) error {
	err := stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Pod",
	})
	if err != nil {
		return err
	}

	pod, err := env.Pod(blackDeath, *cliNamespace, in.Metadata.Name, in.GitCheckout.Repository, in.GitCheckout.Revision, in.Compose.Services)
	if err != nil {
		return fmt.Errorf("failed build pod: %s", err)
	}

	_, err = utils.PodCreate(client, pod)
	if err != nil {
		return fmt.Errorf("failed create pod: %s", err)
	}

	return err
}

// A step for running all our build steps in a single container.
func stepCommands(client *client.Clientset, in *pb.BuildRequest, stream pb.M8S_BuildServer, config *rest.Config, namespace string) error {
	// This is what we will use to communicate back to the CLI client.
	r, w := io.Pipe()

	go func(reader io.Reader, stream pb.M8S_BuildServer) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			err := stream.Send(&pb.BuildResponse{
				Message: scanner.Text(),
			})
			if err != nil {
				fmt.Println("failed to send response:", err)
			}
		}
	}(r, stream)

	// Run the commands inside the pod.
	err := stream.Send(&pb.BuildResponse{
		Message: "Running build steps against K8s Pod",
	})
	if err != nil {
		return err
	}

	for _, step := range in.Exec.Steps {
		err = stream.Send(&pb.BuildResponse{
			Message: fmt.Sprintf("Running command: %s", step),
		})
		if err != nil {
			return err
		}

		err = utils.PodExec(client, config, w, namespace, in.Metadata.Name, in.Exec.Container, step)
		if err != nil {
			return fmt.Errorf("ftep failed: %s", err)
		}
	}

	return nil
}
