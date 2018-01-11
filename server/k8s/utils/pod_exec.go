package utils

import (
	"io"

	remotecommandserver "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const error129 = "command terminated with exit code 129"

// PodExecInput is used for passing params to the PodExec function.
type PodExecInput struct {
	Client    *kubernetes.Clientset
	Config    *rest.Config
	Stdin     bool
	Stdout    bool
	Stderr    bool
	Writer    io.Writer
	Reader    io.Reader
	Namespace string
	Pod       string
	Container string
	Command   []string
}

// PodExec for running commands against a running pod.
// func PodExec(client *kubernetes.Clientset, config *rest.Config, w io.Writer, r io.Reader, namespace, name, container string, command []string) error {
func PodExec(input PodExecInput) error {
	opts := remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Tty:                true,
	}

	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	req := input.Client.Core().RESTClient().Post().Resource("pods").Name(input.Pod).Namespace(input.Namespace).SubResource("exec")
	req.Param("container", input.Container)

	if input.Stdin {
		req.Param("stdin", "true")
		opts.Stdin = input.Reader
	}

	if input.Stdout {
		req.Param("stdout", "true")
		opts.Stdout = input.Writer
	}

	if input.Stderr {
		req.Param("stderr", "true")
		opts.Stderr = input.Writer
	}

	req.Param("tty", "true")
	for _, cmd := range input.Command {
		req.Param("command", cmd)
	}
	url := req.URL()

	exec, err := remotecommand.NewExecutor(input.Config, "POST", url)
	if err != nil {
		return err
	}

	err = exec.Stream(opts)
	// This is not the most ideal way to handle the error since its tied to a specfic string.
	// This appears to be an error from Docker.
	// https://github.com/docker/compose/issues/3379
	// @todo, Compare to containerd backed container runtime.
	if err != nil && err.Error() != error129 {
		return err
	}

	return nil
}
