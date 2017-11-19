package m8s

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func installPVC(client *kubernetes.Clientset, namespace string) error {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m8s-api-tls",
			Namespace: namespace,
			Annotations: map[string]string{
				"volume.beta.kubernetes.io/storage-class": "standard",
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := utils.PersistentVolumeClaimCreate(client, pvc)
	if err != nil {
		return errors.Wrap(err, "failed to install PVC")
	}

	return nil
}
