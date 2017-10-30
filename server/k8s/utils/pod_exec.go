package utils

import (
	"io"

	remotecommandserver "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const error129 = "command terminated with exit code 129"

// PodExec for running commands against a running pod.
func PodExec(client *kubernetes.Clientset, config *rest.Config, w io.Writer, namespace, name, container, step string) error {
	opts := remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdout:             w,
		Stderr:             w,
		Tty:                true,
	}

	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	req := client.Core().RESTClient().Post().Resource("pods").Name(name).Namespace(namespace).SubResource("exec")
	req.Param("container", container)
	req.Param("stdout", "true")
	req.Param("stderr", "true")
	req.Param("tty", "true")
	req.Param("command", "/bin/bash")
	req.Param("command", "-c")
	req.Param("command", step)
	url := req.URL()

	exec, err := remotecommand.NewExecutor(config, "POST", url)
	// This is not the most ideal way to handle the error since its tied to a specfic string.
	// This appears to be an error from Docker.
	// https://github.com/docker/compose/issues/3379
	// @todo, Compare to containerd backed container runtime.
	if err != nil && err.Error() != error129 {
		return err
	}

	return exec.Stream(opts)
}
