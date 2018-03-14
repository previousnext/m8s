package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerate(t *testing.T) {
	have, err := Generate(GenerateParams{
		Namespace: "test",
		Name:      "test",
		Annotations: map[string]string{
			"nick": "rocks",
		},
	})
	assert.Nil(t, err)

	assert.Equal(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
			Annotations: map[string]string{
				"nick": "rocks",
			},
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
				"env": "test",
			},
		},
	}, have)
}
