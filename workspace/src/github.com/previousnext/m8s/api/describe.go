package main

import (
	"fmt"

	"github.com/previousnext/m8s/api/k8s/addons"
	pb "github.com/previousnext/m8s/pb"
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

	ing, err := srv.client.Extensions().Ingresses(*cliNamespace).Get(in.Name, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	svc, err := srv.client.Services(*cliNamespace).Get(addons.SSHName, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	env := &pb.Environment{
		Name:      pod.ObjectMeta.Name,
		Namespace: pod.ObjectMeta.Namespace,
	}

	// Get the list of domains.
	for _, rule := range ing.Spec.Rules {
		env.Domains = append(env.Domains, rule.Host)
	}

	// Get the list of containers.
	for _, container := range pod.Spec.Containers {
		env.Containers = append(env.Containers, &pb.Container{
			Name:  container.Name,
			Image: container.Image,
		})
	}

	// Get the SSH endpoint (load balancer attached to service).
	for _, balancer := range svc.Status.LoadBalancer.Ingress {
		env.SSH = append(env.SSH, fmt.Sprintf("%s:%v", balancer.Hostname, addons.SSHPort))
	}

	resp.Environment = env

	return resp, nil
}
