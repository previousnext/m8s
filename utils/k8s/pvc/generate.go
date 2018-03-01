package pvc

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateParams for generating an PersistentVolumeClaim object.
type GenerateParams struct {
	Namespace    string
	Name         string
	StorageClass string
	Annotations  map[string]string
}

// Generate is used for creating a PersistentVolumeClaim object.
func Generate(params GenerateParams) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   params.Namespace,
			Name:        params.Name,
			Annotations: params.Annotations,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	pvc.ObjectMeta.Annotations["volume.beta.kubernetes.io/storage-class"] = params.StorageClass

	return pvc, nil
}
