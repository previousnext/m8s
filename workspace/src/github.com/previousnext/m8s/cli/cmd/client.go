package cmd

import (
	pb "github.com/previousnext/m8s/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func buildClient(endpoint string) (pb.M8SClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	return pb.NewM8SClient(conn), nil
}
