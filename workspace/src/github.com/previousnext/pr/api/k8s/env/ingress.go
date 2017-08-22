package env

import (
	"fmt"

	"github.com/previousnext/pr/api/k8s/utils"
	"golang.org/x/crypto/bcrypt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateIngress is used for creating the Ingress object.
func CreateIngress(client *client.Clientset, timeout int64, namespace, name, user, pass string, domains []string) error {
	ing, err := Ingress(timeout, namespace, name, domains)
	if err != nil {
		return err
	}

	if user != "" && pass != "" {
		secretName := fmt.Sprintf("%s-auth", name)

		secret, err := SecretBasicAuth(timeout, namespace, secretName, user, pass)
		if err != nil {
			return err
		}

		err = utils.CreateSecret(client, secret)
		if err != nil {
			return err
		}

		// Add basic auth for Traefik.
		ing.ObjectMeta.Annotations["ingress.kubernetes.io/auth-type"] = "basic"
		ing.ObjectMeta.Annotations["ingress.kubernetes.io/auth-secret"] = secretName
	}

	err = utils.CreateIngress(client, ing)
	if err != nil {
		return err
	}

	return nil
}

// SecretBasicAuth is used for generating a "basic auth" secret for our PR environment.
// @todo, Needs a test.
func SecretBasicAuth(timeout int64, namespace, name, user, pass string) (*v1.Secret, error) {
	// Convert our user and pass into a htpasswd file.
	hash, err := htpasswd(pass)
	if err != nil {
		return nil, err
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"skipper.io/black-death": fmt.Sprintf("%v", timeout),
			},
		},
		StringData: map[string]string{
			"auth": fmt.Sprintf("%s:%s", user, hash),
		},
	}, nil
}

// Ingress converts a Docker Compose file into a Kubernetes Ingress object.
func Ingress(timeout int64, namespace, name string, domains []string) (*extensions.Ingress, error) {
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
				"skipper.io/black-death":      fmt.Sprintf("%v", timeout),
			},
		},
	}

	for _, domain := range domains {
		ingress.Spec.Rules = append(ingress.Spec.Rules, extensions.IngressRule{
			Host: domain,
			IngressRuleValue: extensions.IngressRuleValue{
				HTTP: &extensions.HTTPIngressRuleValue{
					Paths: []extensions.HTTPIngressPath{
						{
							Path: "/",
							Backend: extensions.IngressBackend{
								ServiceName: name,
								ServicePort: intstr.FromInt(80),
							},
						},
					},
				},
			},
		})
	}

	return ingress, nil
}

// Helper function for generating a http auth password.
// @todo, Needs a test.
func htpasswd(pass string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordBytes), nil
}
