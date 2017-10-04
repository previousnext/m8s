package ssh_server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/previousnext/m8s/api/k8s/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

const (
	Name = "ssh-server"
	Port = 22
)

// Create will create our ssh-server ingress router.
func Create(client *client.Clientset, namespace, image, version string) error {
	_, err := createSecret(client, namespace)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server secret: %s", err)
	}

	_, err = createDeployment(client, namespace, image, version, 2, 1)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server deployment: %s", err)
	}

	_, err = createService(client, namespace)
	if err != nil {
		return fmt.Errorf("failed deploy ssh server service: %s", err)
	}

	return nil
}

func createSecret(client *client.Clientset, namespace string) (*v1.Secret, error) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      Name,
		},
	}

	key, err := rsa.GenerateKey(rand.Reader, 768)
	if err != nil {
		return secret, err
	}

	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	}

	secret.Data = map[string][]byte{
		"signer": pem.EncodeToMemory(&priv_blk),
	}

	secret, err = client.CoreV1().Secrets(namespace).Create(secret)
	if errors.IsAlreadyExists(err) {
		return secret, nil
	}

	return secret, err
}

func createDeployment(client *client.Clientset, namespace, image, version string, replicas, history int32) (*v1beta1.Deployment, error) {
	dply := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &history,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: Name,
					Labels: map[string]string{
						"name": Name,
					},
					Namespace: namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  Name,
							Image: fmt.Sprintf("%s:%s", image, version),
							Ports: []v1.ContainerPort{
								{
									Name:          "ssh",
									ContainerPort: Port,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "SSH_SIGNER",
									Value: "/etc/signers/signer",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "signer",
									MountPath: "/etc/signers",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "signer",
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: Name,
								},
							},
						},
					},
				},
			},
		},
	}

	return utils.DeploymentCreate(client, dply)
}

func createService(client *client.Clientset, namespace string) (*v1.Service, error) {
	// This automatically deploys a load balancer for this service.
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      Name,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{
					Name:       "ssh",
					Port:       Port,
					TargetPort: intstr.FromInt(Port),
				},
			},
			// This allows us to link this Service to the Pod.
			Selector: map[string]string{
				"name": Name,
			},
		},
	}

	return utils.ServiceCreate(client, svc)
}
