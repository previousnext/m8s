package k8s

import (
	pb "github.com/previousnext/pr/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Service converts a Docker Compose file into a Kubernetes Service object.
func Service(namespace string, in *pb.BuildRequest) (*v1.Service, error) {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      in.Metadata.Name,
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
				"env": in.Metadata.Name,
			},
		},
	}, nil
}
