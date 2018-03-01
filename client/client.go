package client

import (
	"github.com/pkg/errors"

	"github.com/previousnext/m8s/client/k8s"
	"github.com/previousnext/m8s/client/openshift"
	"github.com/previousnext/m8s/client/types"
)

// New returns a new client.
func New(name string, params types.ClientParams) (types.Client, error) {
	if name == openshift.Name {
		return openshift.New(params)
	}

	if name == k8s.Name {
		return k8s.New(params)
	}

	return nil, errors.New("cannot find client")
}
