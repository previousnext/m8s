package env

import (
	"fmt"
	"time"

	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server/k8s/env/htpasswd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// SecretInput provides the Secret function with information to produce a Kubernetes Secret.
type SecretInput struct {
	Namespace   string
	Name        string
	Annotations []*pb.Annotation
	User        string
	Pass        string
	Retention   string
}

// Secret is used for generating a "basic auth" secret for our PR environment.
// @todo, Needs a test.
func Secret(input SecretInput) (*v1.Secret, error) {
	// Convert our user and pass into a htpasswd file.
	hash, err := htpasswd.Hash(input.Pass)
	if err != nil {
		return nil, err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: input.Namespace,
			Name:      input.Name,
			Annotations: map[string]string{
				"author": "m8s",
			},
		},
		StringData: map[string]string{
			"auth": fmt.Sprintf("%s:%s", input.User, hash),
		},
	}

	for _, annotation := range input.Annotations {
		secret.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	if input.Retention != "" {
		unix, err := retentionToUnix(time.Now(), input.Retention)
		if err != nil {
			return secret, err
		}

		secret.ObjectMeta.Annotations["black-death.skpr.io"] = unix
	}

	return secret, nil
}
