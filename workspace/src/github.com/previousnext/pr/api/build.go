package main

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/previousnext/pr/api/k8s/env"
	pb "github.com/previousnext/pr/pb"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

func (srv server) Build(in *pb.BuildRequest, stream pb.PR_BuildServer) error {
	if in.Credentials.Token != *cliToken {
		return fmt.Errorf("token is incorrect")
	}

	// @todo, Verification.

	keep, err := time.ParseDuration(in.Keep)
	if err != nil {
		return fmt.Errorf("failed marshall field 'keep': %s", err)
	}

	// Create a unix timestamp to be used for "Black Death".
	timeout := time.Now().Unix() + keep.Nanoseconds()

	// Step 1 - Create Kubernetes Service object.
	err = stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Service",
	})
	if err != nil {
		return err
	}

	err = env.CreateService(srv.client, timeout, *cliNamespace, in.Metadata.Name)
	if err != nil {
		return fmt.Errorf("failed create service: %s", err)
	}

	// Step 2 - Create Kubernetes Ingress object.
	err = stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Ingress",
	})
	if err != nil {
		return err
	}

	err = env.CreateIngress(srv.client, timeout, *cliNamespace, in.Metadata.Name, in.Metadata.Domains)
	if err != nil {
		return fmt.Errorf("failed create ingress: %s", err)
	}

	// Step 3 - Create Kubernetes Pod object.
	err = stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Pod",
	})
	if err != nil {
		return err
	}

	pod, err := env.CreatePod(srv.client, timeout, *cliNamespace, in.Metadata.Name, in.GitCheckout.Repository, in.GitCheckout.Revision, in.Compose.Services)
	if err != nil {
		return fmt.Errorf("failed create pod: %s", err)
	}

	// This is what we will use to communicate back to the CLI client.
	r, w := io.Pipe()

	go func(reader io.Reader, stream pb.PR_BuildServer) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			err = stream.Send(&pb.BuildResponse{
				Message: scanner.Text(),
			})
			if err != nil {
				fmt.Println("failed to send response:", err)
			}
		}
	}(r, stream)

	// Step 4 - Run the commands inside the pod.
	err = stream.Send(&pb.BuildResponse{
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

		err = runStep(srv.client, srv.config, w, pod, in.Exec.Container, step)
		if err != nil {
			return fmt.Errorf("ftep failed: %s", err)
		}
	}

	return nil
}

// Helper function for running commands against a running pod.
func runStep(client *client.Clientset, config *rest.Config, w io.Writer, pod *v1.Pod, container, step string) error {
	cmd := &api.PodExecOptions{
		Container: container,
		Stdout:    true,
		Stderr:    true,
		Command: []string{
			"/bin/bash",
			"-c",
			step,
		},
	}

	opts := remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdout:             w,
		Stderr:             w,
	}

	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	url := client.Core().RESTClient().Post().Resource("pods").Name(pod.ObjectMeta.Name).Namespace(pod.ObjectMeta.Namespace).SubResource("exec").Param("container", container).VersionedParams(cmd, api.ParameterCodec).URL()

	exec, err := remotecommand.NewExecutor(config, "POST", url)
	if err != nil {
		return err
	}

	return exec.Stream(opts)
}
