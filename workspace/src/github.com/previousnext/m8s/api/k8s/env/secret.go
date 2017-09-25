package env

import (
	"fmt"

	"github.com/previousnext/htpasswd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Secret is used for generating a "basic auth" secret for our PR environment.
// @todo, Needs a test.
func Secret(namespace, name, user, pass string) (*v1.Secret, error) {
	// Convert our user and pass into a htpasswd file.
	hash, err := htpasswd.Hash(pass)
	if err != nil {
		return nil, err
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		StringData: map[string]string{
			"auth": fmt.Sprintf("%s:%s", user, hash),
		},
	}, nil
}
