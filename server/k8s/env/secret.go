package env

import (
	"fmt"
	"time"

	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server/k8s/env/htpasswd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Secret is used for generating a "basic auth" secret for our PR environment.
// @todo, Needs a test.
func Secret(namespace, name string, annotations []*pb.Annotation, user, pass, retention string) (*v1.Secret, error) {
	// Convert our user and pass into a htpasswd file.
	hash, err := htpasswd.Hash(pass)
	if err != nil {
		return nil, err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		StringData: map[string]string{
			"auth": fmt.Sprintf("%s:%s", user, hash),
		},
	}

	for _, annotation := range annotations {
		secret.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	if retention != "" {
		unix, err := retentionToUnix(time.Now(), retention)
		if err != nil {
			return secret, err
		}

		secret.ObjectMeta.Annotations["black-death.skpr.io"] = unix
	}

	return secret, nil
}
