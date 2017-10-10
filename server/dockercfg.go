package server

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	context "golang.org/x/net/context"
)

// DockerCfg returns Docker credentials for pushing and pulling images.
func (srv Server) DockerCfg(ctx context.Context, in *pb.DockerCfgRequest) (*pb.DockerCfgResponse, error) {
	resp := new(pb.DockerCfgResponse)

	if in.Credentials.Token != srv.Token {
		return resp, fmt.Errorf("token is incorrect")
	}

	return &pb.DockerCfgResponse{
		Registry: srv.Docker.Registry,
		Username: srv.Docker.Username,
		Password: srv.Docker.Password,
		Email:    srv.Docker.Email,
		Auth:     srv.Docker.Auth,
	}, nil
}
