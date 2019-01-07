package podutils

import (
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// https://github.com/docker/compose/issues/3379
const error129 = "command terminated with exit code 129"

// ExecParams is used for passing params to the Run function.
type ExecParams struct {
	Client    *kubernetes.Clientset
	Config    *rest.Config
	Stdin     bool
	Stdout    bool
	Stderr    bool
	TTY       bool
	Writer    io.Writer
	Reader    io.Reader
	Namespace string
	Pod       string
	Container string
	Command   []string
}

// Exec a command within a container, within a pod.
func Exec(params ExecParams) error {
	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	req := params.Client.CoreV1().RESTClient().Post().Resource("pods").Name(params.Pod).Namespace(params.Namespace).SubResource("exec")
	req.Param("container", params.Container)
	req.Param("stdin", "true")
	req.Param("stdout", "true")
	req.Param("stderr", "true")
	req.Param("tty", "true")

	opts := remotecommand.StreamOptions{
		Stdin:  params.Reader,
		Stdout: params.Writer,
		Stderr: params.Writer,
		Tty:    true,
	}

	for _, cmd := range params.Command {
		req.Param("command", cmd)
	}
	url := req.URL()

	exec, err := remotecommand.NewSPDYExecutor(params.Config, "POST", url)
	if err != nil {
		return err
	}

	err = exec.Stream(opts)
	if err != nil && err.Error() != error129 {
		return err
	}

	return nil
}
