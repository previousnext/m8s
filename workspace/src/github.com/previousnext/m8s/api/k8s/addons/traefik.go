package addons

import (
	"fmt"

	"github.com/previousnext/m8s/api/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

const (
	TraefikName = "addon-traefik"
	TraefikPort = 80
)

// CreateTraefik will create our Traefik ingress router.
func CreateTraefik(client *client.Clientset, namespace, image, version string, port int32) error {
	var (
		history = int32(1)

		// Deploy this as a HA service, ensuring Ingress will still work.
		replicas = int32(2)
	)

	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TraefikName,
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &history,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": TraefikName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: TraefikName,
					Labels: map[string]string{
						"name": TraefikName,
					},
					Namespace: namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  TraefikName,
							Image: fmt.Sprintf("%s:%s", image, version),
							Args: []string{
								"--web",
								"--kubernetes",
							},
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: TraefikPort,
									HostPort:      port,
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
		return fmt.Errorf("failed deploy traefik: %s", err)
	}

	// This automatically deploys a load balancer for this service.
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      TraefikName,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       TraefikPort,
					TargetPort: intstr.FromInt(TraefikPort),
				},
			},
			// This allows us to link this Service to the Pod.
			Selector: map[string]string{
				"name": TraefikName,
			},
		},
	}

	_, err = utils.ServiceCreate(client, svc)
	if err != nil {
		return fmt.Errorf("failed deploy traefik: %s", err)
	}

	return nil
}
