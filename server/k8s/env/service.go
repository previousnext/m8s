package env

import (
	"time"

	pb "github.com/previousnext/m8s/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// ServiceInput provides the Service function with information to produce a Kubernetes Service.
type ServiceInput struct {
	Namespace   string
	Name        string
	Annotations []*pb.Annotation
	Retention   string
}

// Service converts a Docker Compose file into a Kubernetes Service object.
func Service(input ServiceInput) (*v1.Service, error) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: input.Namespace,
			Name:      input.Name,
			Annotations: map[string]string{
				"author": "m8s",
			},
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None", // We defer this logic to the load balancer.
			Ports: []v1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
				{
					Name: "mailhog",
					Port: 8025,
				},
				{
					Name: "solr",
					Port: 8983,
				},
			},
			// This allows us to Link tihs Service to the Pod.
			Selector: map[string]string{
				"env": input.Name,
			},
		},
	}

	for _, annotation := range input.Annotations {
		svc.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	if input.Retention != "" {
		unix, err := retentionToUnix(time.Now(), input.Retention)
		if err != nil {
			return svc, err
		}

		svc.ObjectMeta.Annotations["black-death.skpr.io"] = unix
	}

	return svc, nil
}
