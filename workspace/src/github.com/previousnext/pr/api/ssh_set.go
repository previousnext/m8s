package main

import (
	"fmt"

	"github.com/previousnext/pr/api/k8s/env"
	pb "github.com/previousnext/pr/pb"
	context "golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

func (srv server) SSHSet(ctx context.Context, in *pb.SSHSetRequest) (*pb.SSHSetResponse, error) {
	resp := new(pb.SSHSetResponse)

	if in.Credentials.Token != *cliToken {
		return resp, fmt.Errorf("token is incorrect")
	}

	if len(in.SSH.PrivateKey) == 0 {
		return resp, fmt.Errorf("private key not found")
	}

	if len(in.SSH.KnownHosts) == 0 {
		return resp, fmt.Errorf("known hosts not found")
	}

	obj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: *cliNamespace,
			Name:      env.SecretSSH,
		},
		Data: map[string][]byte{
			keyPrivateKey: in.SSH.PrivateKey,
			keyKnownHosts: in.SSH.KnownHosts,
		},
		Type: v1.SecretTypeDockercfg,
	}

	_, err := srv.client.Secrets(*cliNamespace).Create(obj)
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
