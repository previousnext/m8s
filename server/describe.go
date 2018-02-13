package server

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Describe returns details about the temporary environment.
func (srv Server) Describe(ctx context.Context, in *pb.DescribeRequest) (*pb.DescribeResponse, error) {
	resp := new(pb.DescribeResponse)

	if in.Credentials.Token != srv.Token {
		return resp, fmt.Errorf("token is incorrect")
	}

	if in.Name == "" {
		return resp, fmt.Errorf("name is incorrect")
	}

	pod, err := srv.client.CoreV1().Pods(srv.Namespace).Get(in.Name, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	ing, err := srv.client.Extensions().Ingresses(srv.Namespace).Get(in.Name, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	svc, err := srv.client.CoreV1().Services(srv.Namespace).Get(srv.SSHService, metav1.GetOptions{})
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
