package env

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// CreatePersistentVolumeClaim is used for creating a new PersistentVolumeClaim object.
func CreatePersistentVolumeClaim(client *client.Clientset, namespace, name, storage string) error {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				// Setting this storage class to "cache" allows system admins to register any type of
				// storage backend for "cache" claims.
				"volume.beta.kubernetes.io/storage-class": "cache",
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteMany,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(storage),
				},
			},
		},
	}

	_, err := client.PersistentVolumeClaims(namespace).Create(pvc)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
