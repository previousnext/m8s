package main

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/previousnext/skpr/packager"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliProject  = kingpin.Flag("project", "Project name").Required().String()
	cliFile     = kingpin.Flag("file", "File to load image configuration").Default("images.yml").String()
	cliRegistry = kingpin.Flag("registry", "Registry to push images to").Default("https://index.docker.io/v1/").String()
	cliUsername = kingpin.Flag("username", "Username for registry authentication").Required().String()
	cliPassword = kingpin.Flag("password", "Password for registry authentication").Required().String()
)

func main() {
	kingpin.Parse()

	var images []packager.Image

	c, err := ioutil.ReadFile(*cliFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(c, &images)
	if err != nil {
		panic(err)
	}

	params := packager.BuildManyParams{
		Config: packager.Config{
			Project: *cliProject,
			Credentials: packager.Credentials{
				Registry: *cliRegistry,
				Username: *cliUsername,
				Password: *cliPassword,
			},
		},
		Images: images,
	}

	err = packager.BuildMany(params, os.Stdout)
	if err != nil {
		panic(err)
	}
}
