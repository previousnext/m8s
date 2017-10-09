package server

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	context "golang.org/x/net/context"
)

func (srv Server) DockerCfg(ctx context.Context, in *pb.DockerCfgRequest) (*pb.DockerCfgResponse, error) {
	resp := new(pb.DockerCfgResponse)

	if in.Credentials.Token != srv.Token {
		return resp, fmt.Errorf("token is incorrect")
	}

	return &pb.DockerCfgResponse{
		Registry: srv.DockerCfgRegistry,
		Username: srv.DockerCfgUsername,
		Password: srv.DockerCfgPassword,
		Email:    srv.DockerCfgEmail,
		Auth:     srv.DockerCfgAuth,
	}, nil
}
