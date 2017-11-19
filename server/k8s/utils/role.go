package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

// RoleCreate is a wrapper which will attempt to create and/or up an Roles.
func RoleCreate(client *kubernetes.Clientset, new *v1beta1.Role) (*v1beta1.Role, error) {
	role, err := client.Rbac().Roles(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.Rbac().Roles(new.ObjectMeta.Namespace).Update(new)
	}

	return role, err
}
