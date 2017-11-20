package m8s

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func installService(client *kubernetes.Clientset, namespace string) error {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m8s-api",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "https",
					Port:       443,
					TargetPort: intstr.FromInt(443),
				},
			},
			Selector: map[string]string{
				"app": "m8s-api",
			},
			Type: v1.ServiceTypeLoadBalancer,
		},
	}

	_, err := utils.ServiceCreate(client, svc)
	if err != nil {
		return errors.Wrap(err, "failed to install Service")
	}

	return nil
}
