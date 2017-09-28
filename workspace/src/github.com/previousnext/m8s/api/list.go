package main

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (srv server) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	resp := new(pb.ListResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	pods, err := srv.client.Pods(*cliNamespace).List(metav1.ListOptions{})
	if err != nil {
		return resp, err
	}

	for _, pod := range pods.Items {
		resp.Environments = append(resp.Environments, pod.ObjectMeta.Name)
	}

	return resp, nil
}
