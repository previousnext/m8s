package env

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

func CreateIngress(client *client.Clientset, namespace, name string, domains []string) error {
	ing, err := Ingress(namespace, name, domains)
	if err != nil {
		return err
	}

	_, err = client.Extensions().Ingresses(namespace).Create(ing)
	if errors.IsAlreadyExists(err) {
		_, err = client.Extensions().Ingresses(namespace).Update(ing)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// Ingress converts a Docker Compose file into a Kubernetes Ingress object.
func Ingress(namespace, name string, domains []string) (*extensions.Ingress, error) {
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
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
