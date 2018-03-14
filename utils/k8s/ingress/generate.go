package ingress

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateParams for generating an Ingress object.
type GenerateParams struct {
	Namespace   string
	Name        string
	Annotations map[string]string
	Domain      string
}

// Generate will generate an Ingress object.
func Generate(params GenerateParams) (*extensionsv1beta1.Ingress, error) {
	ingress := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   params.Namespace,
			Name:        params.Name,
			Annotations: params.Annotations,
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: params.Domain,
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: params.Name,
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ingress, nil
}
