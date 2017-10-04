package main

import (
	"fmt"

	"github.com/previousnext/m8s/api/k8s/addons/ssh-server"
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

	svc, err := srv.client.Services(*cliNamespace).Get(ssh_server.Name, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	resp.Name = pod.ObjectMeta.Name
	resp.Namespace = pod.ObjectMeta.Namespace

	// Get the list of domains.
	for _, rule := range ing.Spec.Rules {
		resp.Domains = append(resp.Domains, rule.Host)
	}

	// Get the list of containers.
	for _, container := range pod.Spec.Containers {
		resp.Containers = append(resp.Containers, container.Name)
	}

	if len(svc.Status.LoadBalancer.Ingress) > 0 {
		resp.SSH = svc.Status.LoadBalancer.Ingress[0].Hostname
	}

	return resp, nil
}
