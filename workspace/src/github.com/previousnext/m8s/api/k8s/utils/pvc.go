package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// PersistentVolumeClaimCreate is used for creating new Persistent Volume Claims.
func PersistentVolumeClaimCreate(client *client.Clientset, new *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	pvc, err := client.PersistentVolumeClaims(new.ObjectMeta.Namespace).Create(new)
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, err
	}

	return pvc, nil
}
