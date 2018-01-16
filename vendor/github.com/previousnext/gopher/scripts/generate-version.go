package main

import (
	"os"

	"github.com/previousnext/gopher/pkg/genversion"
)

const (
	versionFile = "version/version.go"
)

func main() {
	out, err := os.Create(versionFile)
	if err != nil {
		panic(err)
	}

	err = genversion.GenerateVersionFile(out)
	if err != nil {
		panic(err)
	}
}
