package cmd

import (
	pb "github.com/previousnext/m8s/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func buildClient(endpoint string, insecure bool) (pb.M8SClient, error) {
	conn, err := buildConfig(endpoint, insecure)
	if err != nil {
		return nil, err
	}

	return pb.NewM8SClient(conn), nil
}

func buildConfig(endpoint string, insecure bool) (*grpc.ClientConn, error) {
	if insecure {
		return grpc.Dial(endpoint, grpc.WithInsecure())
	}

	return grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
}
