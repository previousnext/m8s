package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateParams for generating an Service object.
type GenerateParams struct {
	Namespace   string
	Name        string
	Annotations map[string]string
}

// Generate creates a Kubernetes Service object.
func Generate(params GenerateParams) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   params.Namespace,
			Name:        params.Name,
			Annotations: params.Annotations,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
			Selector: map[string]string{
				"env": params.Name,
			},
		},
	}

	return svc, nil
}
