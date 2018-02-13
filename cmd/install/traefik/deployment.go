package traefik

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
)

var replicas int32 = 1

// Helper function to install Traefik for Ingress handling.
func installDeployment(client *kubernetes.Clientset, namespace string) error {
	dply := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traefik",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "traefik",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "traefik",
							Image: "traefik:1.3",
							Args: []string{
								"--web",
								"--kubernetes",
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("40Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("40Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := utils.DeploymentCreate(client, dply)
	if err != nil {
		return errors.Wrap(err, "failed to install Deployment")
	}

	return nil
}
