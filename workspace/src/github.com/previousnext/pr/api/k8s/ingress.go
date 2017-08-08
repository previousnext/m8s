package k8s

import (
	pb "github.com/previousnext/pr/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// Ingress converts a Docker Compose file into a Kubernetes Ingress object.
func Ingress(namespace string, in *pb.BuildRequest) (*extensions.Ingress, error) {
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      in.Metadata.Name,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
			},
		},
	}

	for _, domain := range in.Metadata.Domains {
		ingress.Spec.Rules = append(ingress.Spec.Rules, extensions.IngressRule{
			Host: domain,
			IngressRuleValue: extensions.IngressRuleValue{
				HTTP: &extensions.HTTPIngressRuleValue{
					Paths: []extensions.HTTPIngressPath{
						{
							Path: "/",
							Backend: extensions.IngressBackend{
								ServiceName: in.Metadata.Name,
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
