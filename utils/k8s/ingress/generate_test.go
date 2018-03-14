package ingress

import (
	"testing"

	"github.com/stretchr/testify/assert"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestGenerate(t *testing.T) {
	have, err := Generate(GenerateParams{
		Namespace: "test",
		Name:      "test",
		Annotations: map[string]string{
			"nick": "rocks",
		},
		Domain: "www.example.com",
	})
	assert.Nil(t, err)

	assert.Equal(t, &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
			Annotations: map[string]string{
				"nick": "rocks",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: "www.example.com",
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "test",
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}, have)
}
