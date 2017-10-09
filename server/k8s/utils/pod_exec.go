package utils

import (
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	remotecommandserver "k8s.io/apimachinery/pkg/util/remotecommand"
)

// PodExec for running commands against a running pod.
func PodExec(client *kubernetes.Clientset, config *rest.Config, w io.Writer, namespace, name, container, step string) error {
	cmd := &api.PodExecOptions{
		Container: container,
		Stdout:    true,
		Stderr:    true,
		Command: []string{
			"/bin/bash",
			"-c",
			step,
		},
		TTY: true,
	}

	opts := remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdout: w,
		Stderr: w,
		Tty:    true,
	}

	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	url := client.Core().RESTClient().Post().Resource("pods").Name(name).Namespace(namespace).SubResource("exec").Param("container", container).VersionedParams(cmd, api.ParameterCodec).URL()

	exec, err := remotecommand.NewExecutor(config, "POST", url)
	if err != nil {
		return err
	}

	return exec.Stream(opts)
}
