package exec

import (
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// https://github.com/docker/compose/issues/3379
const error129 = "command terminated with exit code 129"

// RunParams is used for passing params to the Run function.
type RunParams struct {
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

// Run a command within a container, within a pod.
func Run(params RunParams) error {
	var opts remotecommand.StreamOptions

	// Use the Kubernetes inbuilt client to build a URL endpoint for running our exec command.
	req := params.Client.CoreV1().RESTClient().Post().Resource("pods").Name(params.Pod).Namespace(params.Namespace).SubResource("exec")
	req.Param("container", params.Container)

	if params.Stdin {
		req.Param("stdin", "true")
		opts.Stdin = params.Reader
	}

	if params.Stdout {
		req.Param("stdout", "true")
		opts.Stdout = params.Writer
	}

	if params.Stderr {
		req.Param("stderr", "true")
		opts.Stderr = params.Writer
	}

	if params.TTY {
		req.Param("tty", "true")
		opts.Tty = true
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
