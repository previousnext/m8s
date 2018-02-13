package env

import (
	"time"

	pb "github.com/previousnext/m8s/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	extensions "k8s.io/api/extensions/v1beta1"
)

// IngressInput provides the Ingress function with information to produce a Kubernetes Ingress.
type IngressInput struct {
	Namespace   string
	Name        string
	Annotations []*pb.Annotation
	Secret      string
	Retention   string
	Domains     []string
}

// Ingress converts a Docker Compose file into a Kubernetes Ingress object.
func Ingress(input IngressInput) (*extensions.Ingress, error) {
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: input.Namespace,
			Name:      input.Name,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
				"author":                      "m8s",
			},
		},
	}

	if input.Secret != "" {
		ingress.ObjectMeta.Annotations["ingress.kubernetes.io/auth-type"] = "basic"
		ingress.ObjectMeta.Annotations["ingress.kubernetes.io/auth-secret"] = input.Secret
	}

	if input.Retention != "" {
		unix, err := retentionToUnix(time.Now(), input.Retention)
		if err != nil {
			return ingress, err
		}

		ingress.ObjectMeta.Annotations["black-death.skpr.io"] = unix
	}

	for _, annotation := range input.Annotations {
		ingress.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	for _, domain := range input.Domains {
		ingress.Spec.Rules = append(ingress.Spec.Rules, extensions.IngressRule{
			Host: domain,
			IngressRuleValue: extensions.IngressRuleValue{
				HTTP: &extensions.HTTPIngressRuleValue{
					Paths: []extensions.HTTPIngressPath{
						{
							Path: "/",
							Backend: extensions.IngressBackend{
								ServiceName: input.Name,
								ServicePort: intstr.FromInt(80),
							},
						},
						{
							Path: "/mailhog",
							Backend: extensions.IngressBackend{
								ServiceName: input.Name,
								ServicePort: intstr.FromInt(8025),
							},
						},
						{
							Path: "/solr",
							Backend: extensions.IngressBackend{
								ServiceName: input.Name,
								ServicePort: intstr.FromInt(8983),
							},
						},
					},
				},
			},
		})
	}

	return ingress, nil
}
