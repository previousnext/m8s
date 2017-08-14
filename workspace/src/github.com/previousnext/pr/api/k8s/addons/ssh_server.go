package addons

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateSSHServer will create our ssh-server ingress router.
// @todo, Look at using a DaemonSet.
func CreateSSHServer(client *client.Clientset, namespace, version string, port int32) error {
	var (
		id      = "addon"
		name    = "ssh-server"
		image   = "previousnext/k8s-ssh-server"
		history = int32(1)

		// Deploy this as a HA service, ensuring SSH Ingress will still work.
		replicas = int32(2)
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
							Ports: []v1.ContainerPort{
								{
									Name:          "ssh",
									ContainerPort: 80,
									HostPort:      port,
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
