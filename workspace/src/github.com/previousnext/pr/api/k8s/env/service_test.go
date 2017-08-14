package env

import (
	"testing"

	pb "github.com/previousnext/pr/pb"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

func TestService(t *testing.T) {
	obj := &pb.BuildRequest{
		Metadata: &pb.Metadata{
			Name: "pr1",
		},
	}

	want := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Port: 80,
				},
			},
			Selector: map[string]string{
				"env": "pr1",
			},
		},
	}

	have, err := Service("test", obj)
	assert.Nil(t, err)
	assert.Equal(t, want, have)
}
