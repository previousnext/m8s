package utils

import (
	"io"

	remotecommandserver "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

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
	if err != nil {
		return err
	}

	return exec.Stream(opts)
}
