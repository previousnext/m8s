package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

func TestService(t *testing.T) {
	want := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
			Annotations: map[string]string{
				"skipper.io/black-death": "123456789",
			},
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

	have, err := Service(123456789, "test", "pr1")
	assert.Nil(t, err)
	assert.Equal(t, want, have)
}
