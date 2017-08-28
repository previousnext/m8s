package compose

import (
	"io/ioutil"

	pb "github.com/previousnext/m8s/pb"
	"gopkg.in/yaml.v2"
)

// DockerCompose is an object which encapsulates a Docker Compose file.
type DockerCompose struct {
	Services map[string]DockerComposeService
}

// DockerComposeService a service declared in a Docker Compose file.
type DockerComposeService struct {
	Image       string   `yaml:"image"`
	Build       string   `yaml:"build"`
	Volumes     []string `yaml:"volumes"`
	Ports       []string `yaml:"ports"`
	Environment []string `yaml:"environment"`
	Tmpfs       []string `yaml:"tmpfs"`
}

// GRPC is used for marshalling a Docker Compose file into a PR GRPC object.
func (dc DockerCompose) GRPC() *pb.Compose {
	resp := new(pb.Compose)

	for name, service := range dc.Services {
		resp.Services = append(resp.Services, &pb.ComposeService{
			Name:        name,
			Image:       service.Image,
			Volumes:     service.Volumes,
			Ports:       service.Ports,
			Environment: service.Environment,
			Tmpfs:       service.Tmpfs,
		})
	}

	return resp
}

// Load the Docker Compose file.
func Load(path string) (DockerCompose, error) {
	var dc DockerCompose

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return dc, err
	}

	err = yaml.Unmarshal(file, &dc)
	if err != nil {
		return dc, err
	}

	return dc, nil
}
