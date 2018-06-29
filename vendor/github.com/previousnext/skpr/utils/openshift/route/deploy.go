package route

import (
	routev1 "github.com/openshift/api/route/v1"
	routev1client "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Deploy will create the Service if not present.
func Deploy(client *routev1client.RouteV1Client, route *routev1.Route) error {
	_, err := client.Routes(route.ObjectMeta.Namespace).Create(route)
	if kerrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
