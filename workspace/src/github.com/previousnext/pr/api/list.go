package main

import (
	"fmt"

	pb "github.com/previousnext/pr/pb"
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
	}

	return resp, nil
}
