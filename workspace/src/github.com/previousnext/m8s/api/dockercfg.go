package main

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	context "golang.org/x/net/context"
)

func (srv server) DockerCfg(ctx context.Context, in *pb.DockerCfgRequest) (*pb.DockerCfgResponse, error) {
	resp := new(pb.DockerCfgResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	return &pb.DockerCfgResponse{
		Registry: *cliDockerCfgRegistry,
		Username: *cliDockerCfgUsername,
		Password: *cliDockerCfgPassword,
		Email:    *cliDockerCfgEmail,
		Auth:     *cliDockerCfgAuth,
	}, nil
}
