package main

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/previousnext/pr/api/k8s"
	pb "github.com/previousnext/pr/pb"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

func (srv server) Build(in *pb.BuildRequest, stream pb.PR_BuildServer) error {
	if in.Credentials.Token != *cliToken {
		return fmt.Errorf("token is incorrect")
	}

	err := stream.Send(&pb.BuildResponse{
		Message: "Mashalling request into K8s objects",
	})
	if err != nil {
		return err
	}

	service, err := k8s.Service(*cliNamespace, in)
	if err != nil {
		return fmt.Errorf("failed to build K8s Service object: %s", err)
	}

	ingress, err := k8s.Ingress(*cliNamespace, in)
	if err != nil {
		return fmt.Errorf("failed to build K8s Ingress object: %s", err)
	}

	pod, err := k8s.Pod(*cliNamespace, in)
	if err != nil {
		return fmt.Errorf("failed to build K8s Pod object: %s", err)
	}

	// Step 1 - Create Kubernetes Service object.
	err = stream.Send(&pb.BuildResponse{
		Message: "Creating K8s Service",
	})
	if err != nil {
		return err
	}

	err = createService(srv.client, service)
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

	err = createIngress(srv.client, ingress)
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

	err = createPod(srv.client, pod)
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

// Helper function for spinning up a new service.
func createService(client *client.Clientset, service *v1.Service) error {
	_, err := client.Services(service.ObjectMeta.Namespace).Create(service)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// Helper function for spinning up a new ingress.
func createIngress(client *client.Clientset, ingress *extensions.Ingress) error {
	_, err := client.Extensions().Ingresses(ingress.ObjectMeta.Namespace).Create(ingress)
	if errors.IsAlreadyExists(err) {
		_, err = client.Extensions().Ingresses(ingress.ObjectMeta.Namespace).Update(ingress)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// Helper function for spinning up a new pod.
func createPod(client *client.Clientset, pod *v1.Pod) error {
	_, err := client.Pods(pod.ObjectMeta.Namespace).Create(pod)
	if errors.IsAlreadyExists(err) {
		// This will tell Kubernetes that we want this pod to be deleted immediately.
		now := int64(0)

		// Delete the Pod.
		err = client.Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{
			GracePeriodSeconds: &now,
		})
		if err != nil {
			return err
		}

		// Create the new pod.
		_, err = client.Pods(pod.ObjectMeta.Namespace).Create(pod)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Wait for the pod to become available.
	limiter := time.Tick(time.Second / 10)

	for {
		pod, err = client.Pods(pod.ObjectMeta.Namespace).Get(pod.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pod.Status.Phase == v1.PodRunning {
			break
		}

		<-limiter
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
