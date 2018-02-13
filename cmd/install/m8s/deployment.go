package m8s

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

func installDeployment(client *kubernetes.Clientset, namespace, token, letsEncryptDomain, letsEncryptEmail string) error {
	dply := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m8s-api",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "m8s-api",
					},
					Annotations: map[string]string{
						"prometheus.io/port":   "9000",
						"prometheus.io/scrape": "true",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "api",
							Image: "previousnext/m8s:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "M8S_TOKEN",
									Value: token,
								},
								{
									Name:  "M8S_NAMESPACE",
									Value: namespace,
								},
								{
									Name:  "M8S_LETS_ENCRYPT_DOMAIN",
									Value: letsEncryptDomain,
								},
								{
									Name:  "M8S_LETS_ENCRYPT_EMAIL",
									Value: letsEncryptEmail,
								},
								{
									// Leaving this blank for now......
									Name:  "M8S_CACHE_DIRS",
									Value: "",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("40Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
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
