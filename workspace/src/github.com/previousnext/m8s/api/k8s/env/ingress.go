package env

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// Ingress converts a Docker Compose file into a Kubernetes Ingress object.
func Ingress(namespace, name, secret string, domains []string) (*extensions.Ingress, error) {
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
			},
		},
	}

	if secret != "" {
		ingress.ObjectMeta.Annotations["ingress.kubernetes.io/auth-type"] = "basic"
		ingress.ObjectMeta.Annotations["ingress.kubernetes.io/auth-secret"] = secret
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
						{
							Path: "/mailhog",
							Backend: extensions.IngressBackend{
								ServiceName: name,
								ServicePort: intstr.FromInt(8025),
							},
						},
						{
							Path: "/solr",
							Backend: extensions.IngressBackend{
								ServiceName: name,
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
