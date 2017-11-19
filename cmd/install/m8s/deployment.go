package m8s

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

var replicas int32 = 1

func installDeployment(client *kubernetes.Clientset, namespace, token, letsEncryptDomain, letsEncryptEmail string) error {
	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m8s-api",
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "m8s-api",
					},
					Annotations: map[string]string{
						"prometheus.io/port":   "9000",
						"prometheus.io/scrape": "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "api",
							Image: "previousnext/m8s:latest",
							Env: []v1.EnvVar{
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
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("40Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("40Mi"),
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
