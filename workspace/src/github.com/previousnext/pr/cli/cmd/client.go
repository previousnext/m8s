package cmd

import (
	pb "github.com/previousnext/pr/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func buildClient(endpoint string) (pb.PRClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	return pb.NewPRClient(conn), nil
}
