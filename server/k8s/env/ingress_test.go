package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

func TestIngress(t *testing.T) {
	want := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":       "traefik",
				"ingress.kubernetes.io/auth-type":   "basic",
				"ingress.kubernetes.io/auth-secret": "pr1",
			},
		},
		Spec: extensions.IngressSpec{
			Rules: []extensions.IngressRule{
				{
					Host: "pr1.example.com",
					IngressRuleValue: extensions.IngressRuleValue{
						HTTP: &extensions.HTTPIngressRuleValue{
							Paths: []extensions.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(80),
									},
								},
								{
									Path: "/mailhog",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(8025),
									},
								},
								{
									Path: "/solr",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(8983),
									},
								},
							},
						},
					},
				},
				{
					Host: "pr1.example2.com",
					IngressRuleValue: extensions.IngressRuleValue{
						HTTP: &extensions.HTTPIngressRuleValue{
							Paths: []extensions.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(80),
									},
								},
								{
									Path: "/mailhog",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(8025),
									},
								},
								{
									Path: "/solr",
									Backend: extensions.IngressBackend{
										ServiceName: "pr1",
										ServicePort: intstr.FromInt(8983),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	have, err := Ingress("test", "pr1", "pr1", []string{"pr1.example.com", "pr1.example2.com"})
	assert.Nil(t, err)
	assert.Equal(t, want, have)
}
