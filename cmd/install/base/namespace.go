package base

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func installNamespace(client *kubernetes.Clientset, namespace string) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m8s",
			Namespace: namespace,
		},
	}

	_, err := utils.NamespaceCreate(client, ns)
	if err != nil {
		return errors.Wrap(err, "failed to install Role")
	}

	return nil
}
