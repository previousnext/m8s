package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	rbacv1 "k8s.io/api/rbac/v1"
)

// RoleCreate is a wrapper which will attempt to create and/or up an Roles.
func RoleCreate(client *kubernetes.Clientset, new *rbacv1.Role) (*rbacv1.Role, error) {
	role, err := client.RbacV1().Roles(new.ObjectMeta.Namespace).Create(new)
	if errors.IsAlreadyExists(err) {
		return client.RbacV1().Roles(new.ObjectMeta.Namespace).Update(new)
	}

	return role, err
}
