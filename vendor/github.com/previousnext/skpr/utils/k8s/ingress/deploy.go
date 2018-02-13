package ingress

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"github.com/pkg/errors"
)

// Deploy will create the Ingress and fallback to updating if it exists.
func Deploy(client *kubernetes.Clientset, ingress *extensionsv1beta1.Ingress) error {
	_, err := client.ExtensionsV1beta1().Ingresses(ingress.ObjectMeta.Namespace).Create(ingress)
	if kerrors.IsAlreadyExists(err) {
		_, err := client.ExtensionsV1beta1().Ingresses(ingress.ObjectMeta.Namespace).Update(ingress)
		if err != nil {
			return errors.Wrap(err, "failed to update")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
