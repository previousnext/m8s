package packager

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"

	colorprefix "github.com/previousnext/skpr/utils/color/prefix"
)

// BuildManyParams is used for passing params to the Run function.
type BuildManyParams struct {
	Config Config  `json:"config" yaml:"config"`
	Images []Image `json:"images" yaml:"images"`
}

// BuildMany starts the build packaging pipeline, running multiple image builders at once.
func BuildMany(params BuildManyParams, w io.Writer) error {
	if len(params.Images) == 0 {
		return errors.New("no images were provided")
	}

	for _, image := range params.Images {
		err := Build(params.Config, image, w)
		if err != nil {
			return err
		}
	}

	return nil
}

// Build will package the image and push to the Docker Registry.
func Build(config Config, image Image, w io.Writer) error {
	err := config.Validate()
	if err != nil {
		return errors.Wrap(err, "failed to validate build config")
	}

	err = image.Validate()
	if err != nil {
		return errors.Wrap(err, "failed to image config")
	}

	name := fmt.Sprintf("%s:%s", image.Repository, image.Tag)

	cred := docker.AuthConfiguration{
		Username: config.Credentials.Username,
		Password: config.Credentials.Password,
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()

	go func(prefix string, reader io.Reader) {
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			fmt.Fprintf(w, "%s %s\n", colorprefix.Wrap(prefix), scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(w, "%s %s\n", colorprefix.Wrap(prefix), err)
		}
	}(image.Name, reader)

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         name,
		Dockerfile:   image.Dockerfile,
		Pull:         true,
		OutputStream: writer,
		ContextDir:   image.Context,
		AuthConfigs: docker.AuthConfigurations{
			Configs: map[string]docker.AuthConfiguration{
				config.Credentials.Registry: cred,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to build image")
	}
	defer client.RemoveImage(name)

	err = client.TagImage(name, docker.TagImageOptions{
		Repo:  image.Repository,
		Tag:   image.Tag,
		Force: true,
	})
	if err != nil {
		return errors.Wrap(err, "failed to tag image")
	}

	if !config.Push {
		return nil
	}

	fmt.Fprintf(w, "%s Pushing image to repository\n", colorprefix.Wrap(image.Name))

	err = client.PushImage(docker.PushImageOptions{
		Name: image.Repository,
		Tag:  image.Tag,
	}, cred)
	if err != nil {
		return errors.Wrap(err, "failed to push image")
	}

	fmt.Fprintf(w, "%s Finished!\n", colorprefix.Wrap(image.Name))

	return nil
}
