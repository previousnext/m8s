package env

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreateService is used for creating a new Service object.
func CreateService(client *client.Clientset, timeout int64, namespace, name string) error {
	svc, err := Service(timeout, namespace, name)
	if err != nil {
		return err
	}

	_, err = client.Services(namespace).Create(svc)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// Service converts a Docker Compose file into a Kubernetes Service object.
func Service(timeout int64, namespace, name string) (*v1.Service, error) {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"skipper.io/black-death": fmt.Sprintf("%v", timeout),
			},
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None", // We defer this logic to the load balancer.
			Ports: []v1.ServicePort{
				{
					Port: 80,
				},
			},
			// This allows us to Link tihs Service to the Pod.
			Selector: map[string]string{
				"env": name,
			},
		},
	}, nil
}
