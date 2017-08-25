package main

import (
	"fmt"

	pb "github.com/previousnext/pr/pb"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (srv server) Describe(ctx context.Context, in *pb.DescribeRequest) (*pb.DescribeResponse, error) {
	resp := new(pb.DescribeResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	if in.Name == "" {
		return resp, fmt.Errorf("name is incorrect")
	}

	pod, err := srv.client.Pods(*cliNamespace).Get(in.Name, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	env := &pb.Environment{
		Name:      pod.ObjectMeta.Name,
		Namespace: pod.ObjectMeta.Namespace,
	}

	for _, container := range pod.Spec.Containers {
		env.Containers = append(env.Containers, &pb.Container{
			Name:  container.Name,
			Image: container.Image,
		})
	}

	return resp, nil
}
