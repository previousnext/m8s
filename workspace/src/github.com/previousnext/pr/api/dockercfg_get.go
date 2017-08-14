package main

import (
	"encoding/json"
	"fmt"

	"github.com/previousnext/pr/api/k8s/env"
	pb "github.com/previousnext/pr/pb"
	context "golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (srv server) DockerCfgGet(ctx context.Context, in *pb.DockerCfgGetRequest) (*pb.DockerCfgGetResponse, error) {
	resp := new(pb.DockerCfgGetResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	secret, err := srv.client.Secrets(*cliNamespace).Get(env.SecretDockerCfg, metav1.GetOptions{})
	if err != nil {
		return resp, err
	}

	cfg, err := getDockerConfig(secret.Data[keyDockerCfg])
	if err != nil {
		return resp, err
	}

	resp.DockerCfg = cfg

	return resp, nil
}

func getDockerConfig(data []byte) (*pb.DockerCfg, error) {
	var dockercfg map[string]DockerConfig

	err := json.Unmarshal(data, &dockercfg)
	if err != nil {
		return nil, err
	}

	// This isn't the best way to handle this, but we are assuming that we will only
	// ever store 1 registry configuration.
	for registry, cfg := range dockercfg {
		return &pb.DockerCfg{
			Registry: registry,
			Username: cfg.Username,
			Password: cfg.Password,
			Email:    cfg.Email,
			Auth:     cfg.Auth,
		}, nil
	}

	return nil, fmt.Errorf("cannot find Docker configuration")
}
