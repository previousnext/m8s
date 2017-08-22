package addons

import (
	"fmt"

	"github.com/previousnext/pr/api/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

const bdName = "black-death"

// CreateBlackDeath is used for creating the "Black Death" addon.
func CreateBlackDeath(client *client.Clientset, namespace, image, version string) error {
	var (
		id       = "addon"
		history  = int32(1)
		replicas = int32(1)
	)

	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", id, bdName),
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &history,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					id: bdName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", id, bdName),
					Labels: map[string]string{
						id: bdName,
					},
					Namespace: namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  bdName,
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

	return utils.CreateDeployment(client, dply)
}
