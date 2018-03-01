package version

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/gosuri/uitable"
)

// PrintParams are passed to the Print function.
type PrintParams struct {
	Version string
	Commit string
	OS string
	Arch string
}

// Print out the version information.
func Print(w io.Writer, params PrintParams) error {
	if params.Version == "" {
		return errors.New("version not found")
	}

	if params.Commit == "" {
		return errors.New("commit not found")
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("Version:", params.Version)
	table.AddRow("Commit:", params.Commit)
	table.AddRow("OS:", params.OS)
	table.AddRow("Arch:", params.Arch)

	fmt.Fprintln(w, table)

	return nil
}
