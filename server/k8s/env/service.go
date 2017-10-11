package env

import (
	pb "github.com/previousnext/m8s/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Service converts a Docker Compose file into a Kubernetes Service object.
func Service(namespace, name string, annotations []*pb.Annotation) *v1.Service {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Annotations: make(map[string]string),
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
				"env": name,
			},
		},
	}

	for _, annotation := range annotations {
		svc.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	return svc
}
