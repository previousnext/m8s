package main

import (
	pb "github.com/previousnext/m8s/pb"
)

func authProvided(in *pb.CreateRequest) bool {
	if in.Metadata.BasicAuth == nil {
		return false
	}

	if in.Metadata.BasicAuth.User == "" {
		return false
	}

	if in.Metadata.BasicAuth.Pass == "" {
		return false
	}

	return true
}
