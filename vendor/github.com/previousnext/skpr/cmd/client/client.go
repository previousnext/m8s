package client

import "io"

// Client interface for allowing multiple backends.
type Client interface {
	Deploy(io.Writer, DeployParams) error
	Package(io.Writer, PackageParams) error
	Version(io.Writer, VersionParams) error
}

// DeployParams are passed to the Deploy function.
type DeployParams struct {
	Environment string
	Version     string
	DryRun      bool
	Wait        bool
}

// PackageParams are passed to the Package function.
type PackageParams struct {
	Version   string
	Directory string
}

// VersionParams are passed to the Version function.
type VersionParams struct{}
