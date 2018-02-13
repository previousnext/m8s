package utils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

// PersistentVolumeClaimCreate is used for creating new Persistent Volume Claims.
func PersistentVolumeClaimCreate(client *kubernetes.Clientset, new *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	pvc, err := client.CoreV1().PersistentVolumeClaims(new.ObjectMeta.Namespace).Create(new)
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, err
	}

	return pvc, nil
}
