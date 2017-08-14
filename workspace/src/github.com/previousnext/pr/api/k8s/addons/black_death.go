package addons

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

func CreateBlackDeath(client *client.Clientset, namespace, version string) error {
	var (
		id       = "addon"
		name     = "black-death"
		image    = "previousnext/k8s-black-death"
		history  = int32(1)
		replicas = int32(1)
	)

	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", id, name),
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &history,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					id: name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", id, name),
					Labels: map[string]string{
						id: name,
					},
					Namespace: namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  name,
							Image: fmt.Sprintf("%s:%s", image, version),
							Env: []v1.EnvVar{
								{
									Name:  "NAMESPACE",
									Value: namespace,
								},
							},
						},
					},
				},
			},
		},
	}

	return createDeployment(client, dply)
}
