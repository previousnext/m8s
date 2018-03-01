package types

import (
	"io"
	"github.com/previousnext/m8s/config"
	"github.com/previousnext/compose"
)

type Client interface {
	Build(io.Writer, BuildParams) error
	Step(io.Writer, StepParams) error
}

type ClientParams struct {
	Master string
	KubeConfig string
}

type BuildParams struct {
	Name string
	Domain string
	Annotations map[string]string
	Repository string
	Revision string
	Config config.Config
	DockerCompose compose.DockerCompose
}

type StepParams struct {
	Namespace string
	Name      string
	Container string
	Command   string
}