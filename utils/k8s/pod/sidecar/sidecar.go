package sidecar

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// DefaultPort for this sidecar.
const DefaultPort int = 8888

// GenerateParams for generating an Sidecar Container object.
type GenerateParams struct {
	Port int
	User string
	Pass string
}

// Generate a Sidecar Container object.
func Generate(params GenerateParams) (corev1.Container, error) {
	if params.Port == 0 {
		params.Port = DefaultPort
	}

	container := corev1.Container{
		Name:  "m8s",
		Image: "previousnext/m8s-router:latest",
		Env: []corev1.EnvVar{
			{
				Name:  "HTTP_PORT",
				Value: strconv.Itoa(params.Port),
			},
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("50m"),
				corev1.ResourceMemory: resource.MustParse("40Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("50m"),
				corev1.ResourceMemory: resource.MustParse("40Mi"),
			},
		},
	}

	if params.User != "" {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "BASIC_AUTH_USER",
			Value: params.User,
		})
	}

	if params.Pass != "" {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "BASIC_AUTH_PASS",
			Value: params.Pass,
		})
	}

	return container, nil
}
