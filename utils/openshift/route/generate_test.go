package route

import (
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestGenerate(t *testing.T) {
	have, err := Generate(GenerateParams{
		Namespace: "test",
		Name:      "test",
		Annotations: map[string]string{
			"nick": "rocks",
		},
		Domain: "www.example.com",
		Port:   80,
	})
	assert.Nil(t, err)

	assert.Equal(t, &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
			Annotations: map[string]string{
				"nick": "rocks",
			},
		},
		Spec: routev1.RouteSpec{
			Host: "www.example.com",
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(80),
			},
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: "test",
			},
		},
	}, have)
}
