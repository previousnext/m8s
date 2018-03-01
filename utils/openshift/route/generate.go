package route

import (
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateParams for generating an Route object.
type GenerateParams struct {
	Namespace   string
	Name        string
	Annotations map[string]string
	Domain      string
}

// Generate will generate an Route object.
func Generate(params GenerateParams) (*routev1.Route, error) {
	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   params.Namespace,
			Name:        params.Name,
			Annotations: params.Annotations,
		},
		Spec: routev1.RouteSpec{
			// @todo, Support mailhog + solr.
			Host: params.Domain,
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: params.Name,
			},
		},
	}

	return route, nil
}
