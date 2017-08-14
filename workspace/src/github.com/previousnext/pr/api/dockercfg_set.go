package main

import (
	"encoding/json"
	"fmt"

	"github.com/previousnext/pr/api/k8s/env"
	pb "github.com/previousnext/pr/pb"
	context "golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

func (srv server) DockerCfgSet(ctx context.Context, in *pb.DockerCfgSetRequest) (*pb.DockerCfgSetResponse, error) {
	resp := new(pb.DockerCfgSetResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	if in.DockerCfg.Registry == "" {
		return resp, fmt.Errorf("registry not found")
	}

	if in.DockerCfg.Username == "" {
		return resp, fmt.Errorf("username not found")
	}

	if in.DockerCfg.Password == "" {
		return resp, fmt.Errorf("password not found")
	}

	if in.DockerCfg.Email == "" {
		return resp, fmt.Errorf("email not found")
	}

	if in.DockerCfg.Auth == "" {
		return resp, fmt.Errorf("auth token not found")
	}

	auths := map[string]DockerConfig{
		in.DockerCfg.Registry: {
			Username: in.DockerCfg.Username,
			Password: in.DockerCfg.Password,
			Email:    in.DockerCfg.Email,
			Auth:     in.DockerCfg.Auth,
		},
	}

	dockerconfig, err := json.Marshal(auths)
	if err != nil {
		return resp, err
	}

	obj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: *cliNamespace,
			Name:      env.SecretDockerCfg,
		},
		Data: map[string][]byte{
			keyDockerCfg: dockerconfig,
		},
		Type: v1.SecretTypeDockercfg,
	}

	_, err = srv.client.Secrets(*cliNamespace).Create(obj)
	if errors.IsAlreadyExists(err) {
		_, err = srv.client.Secrets(*cliNamespace).Update(obj)
		if err != nil {
			return resp, err
		}
	} else if err != nil {
		return resp, err
	}

	return resp, nil
}
