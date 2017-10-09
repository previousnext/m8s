package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// PersistentVolumeClaimCreate is used for creating new Persistent Volume Claims.
func PersistentVolumeClaimCreate(client *kubernetes.Clientset, new *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	pvc, err := client.PersistentVolumeClaims(new.ObjectMeta.Namespace).Create(new)
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, err
	}

	return pvc, nil
}
