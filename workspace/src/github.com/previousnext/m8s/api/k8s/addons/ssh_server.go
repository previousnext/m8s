package addons

import (
	"fmt"

	"github.com/previousnext/m8s/api/k8s/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

const (
	SSHName = "ssh-server"
	SSHPort = 22
)

// CreateSSHServer will create our ssh-server ingress router.
func CreateSSHServer(client *client.Clientset, namespace, image, version, storage string) error {
	var (
		history = int32(1)

		// Deploy this as a HA service, ensuring SSH Ingress will still work.
		replicas = int32(2)
	)

	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SSHName,
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &history,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": SSHName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: SSHName,
					Labels: map[string]string{
						"name": SSHName,
					},
					Namespace: namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  SSHName,
							Image: fmt.Sprintf("%s:%s", image, version),
							Ports: []v1.ContainerPort{
								{
									Name:          "ssh",
									ContainerPort: SSHPort,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "SSH_SIGNER",
									Value: "/etc/ssh/signer",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "ssh",
									MountPath: "/etc/ssh",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "ssh",
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: SSHName,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := utils.DeploymentCreate(client, dply)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server deployment: %s", err)
	}

	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      SSHName,
			Annotations: map[string]string{
				// Setting this storage class to "cache" allows system admins to register any type of
				// storage backend for "cache" claims.
				"volume.beta.kubernetes.io/storage-class": "ssh",
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

	_, err = utils.PersistentVolumeClaimCreate(client, pvc)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server pvc: %s", err)
	}

	// This automatically deploys a load balancer for this service.
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      SSHName,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{
					Name:       "ssh",
					Port:       SSHPort,
					TargetPort: intstr.FromInt(SSHPort),
				},
			},
			// This allows us to link this Service to the Pod.
			Selector: map[string]string{
				"name": SSHName,
			},
		},
	}

	_, err = utils.ServiceCreate(client, svc)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server service: %s", err)
	}

	return nil
}
