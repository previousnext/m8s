package traefik

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
)

func installService(client *kubernetes.Clientset, namespace string) error {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traefik",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app": "traefik",
			},
			Type: v1.ServiceTypeLoadBalancer,
		},
	}

	_, err := utils.ServiceCreate(client, svc)
	if err != nil {
		return errors.Wrap(err, "failed to install Service")
	}

	// @todo, Wait for service to get an IP.

	return nil
}
