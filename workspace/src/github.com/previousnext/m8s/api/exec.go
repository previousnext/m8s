package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/previousnext/m8s/api/k8s/utils"
	pb "github.com/previousnext/m8s/pb"
)

func (srv server) Exec(in *pb.ExecRequest, stream pb.M8S_ExecServer) error {
	if in.Credentials.Token != *cliToken {
		return fmt.Errorf("token is incorrect")
	}

	if in.Name == "" {
		return fmt.Errorf("pod name was not provided")
	}

	if in.Container == "" {
		return fmt.Errorf("exec container was not provided")
	}

	if in.Command == "" {
		return fmt.Errorf("exec container was not provided")
	}

	// This is what we will use to communicate back to the CLI client.
	r, w := io.Pipe()

	go func(reader io.Reader, stream pb.M8S_ExecServer) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			err := stream.Send(&pb.ExecResponse{
				Message: scanner.Text(),
			})
			if err != nil {
				fmt.Println("failed to send response:", err)
			}
		}
	}(r, stream)

	err := stream.Send(&pb.ExecResponse{
		Message: fmt.Sprintf("Running command: %s", in.Command),
	})
	if err != nil {
		return err
	}

	err = utils.PodExec(srv.client, srv.config, w, *cliNamespace, in.Name, in.Container, in.Command)
	if err != nil {
		return fmt.Errorf("command failed: %s", err)
	}

	return nil
}
