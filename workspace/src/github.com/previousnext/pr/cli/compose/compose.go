package compose

import (
	"io/ioutil"

	pb "github.com/previousnext/pr/pb"
	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
	Services map[string]DockerComposeService
}

type DockerComposeService struct {
	Image       string   `yaml:"image"`
	Build       string   `yaml:"build"`
	Volumes     []string `yaml:"volumes"`
	Ports       []string `yaml:"ports"`
	Environment []string `yaml:"environment"`
}

func (dc DockerCompose) Proto() *pb.Compose {
	resp := new(pb.Compose)

	for name, service := range dc.Services {
		resp.Services = append(resp.Services, &pb.ComposeService{
			Name:        name,
			Image:       service.Image,
			Volumes:     service.Volumes,
			Ports:       service.Ports,
			Environment: service.Environment,
		})
	}

	return resp
}

// Helper function to load the Docker Compose file.
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
