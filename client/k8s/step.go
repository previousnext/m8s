package k8s

import (
	"io"

	"github.com/previousnext/m8s/client/types"
	"github.com/previousnext/skpr/utils/k8s/pods/exec"
)

func (c Client) Step(w io.Writer, params types.StepParams) error {
	return exec.Run(exec.RunParams{
		Client:    c.client,
		Config:    c.config,
		Stdout:    true,
		Stderr:    true,
		Writer:    w,
		Namespace: params.Namespace,
		Pod:       params.Name,
		Container: params.Container,
		Command: []string{
			"/bin/bash",
			"-c",
			params.Command,
		},
	})
}