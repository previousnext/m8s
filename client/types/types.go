package types

import (
	"github.com/previousnext/compose"
	"github.com/previousnext/m8s/config"
	"io"
)

// Client workflow.
type Client interface {
	Build(io.Writer, BuildParams) error
	Step(io.Writer, StepParams) error
}

// ClientParams are passed to the New() function.
type ClientParams struct {
	Master     string
	KubeConfig string
}

// BuildParams are params for building an environment.
type BuildParams struct {
	Name          string
	Domain        string
	Annotations   map[string]string
	Repository    string
	Revision      string
	Config        config.Config
	DockerCompose compose.DockerCompose
}

// StepParams are params for stepping through a build.
type StepParams struct {
	Namespace string
	Name      string
	Container string
	Command   string
}
