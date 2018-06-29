package pod

import (
	"testing"

	"github.com/previousnext/compose"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerate(t *testing.T) {
	perm := int32(256)

	have, err := Generate(GenerateParams{
		Namespace: "test",
		Name:      "test",
		Annotations: map[string]string{
			"nick": "rocks",
		},
		Repository: "https://github.com/example/repo.git",
		Revision:   "abcdef123456",
		Services: map[string]compose.Service{
			"app": {
				Image: "test/app:0.0.1",
				Volumes: []string{
					".:/data",
				},
				Entrypoint: []string{
					"echo 1",
					"echo 2",
				},
				Ports: []string{
					"80:80",
				},
				Environment: []string{
					"FOO=bar",
				},
				CapAdd: []string{
					"ALLOFIT",
				},
				Tmpfs: []string{
					"/tmp",
				},
				Deploy: compose.ServiceDeploy{
					Resources: compose.ServiceDeployResources{
						Limits: compose.ServiceDeployResource{
							CPUs:   "10m",
							Memory: "128m",
						},
						Reservations: compose.ServiceDeployResource{
							CPUs:   "10m",
							Memory: "128m",
						},
					},
				},
			},
		},
		Caches: []string{
			"/root/.composer",
		},
		SecretDockerCfg: "test-dockercfg",
		SecretSSH:       "test-ssh",
	})
	assert.Nil(t, err)

	assert.Equal(t, metav1.ObjectMeta{
		Namespace: "test",
		Name:      "test",
		Labels: map[string]string{
			"env": "test",
		},
		Annotations: map[string]string{
			"nick": "rocks",
		},
	}, have.ObjectMeta)

	assert.Equal(t, []corev1.Container{
		{
			Name:            "app",
			Image:           "test/app:0.0.1",
			ImagePullPolicy: corev1.PullAlways,
			Command: []string{
				"echo 1",
				"echo 2",
			},
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: int32(80),
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "FOO",
					Value: "bar",
				},
			},
			SecurityContext: &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{
						"ALLOFIT",
					},
				},
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("10m"),
					corev1.ResourceMemory: resource.MustParse("128m"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("10m"),
					corev1.ResourceMemory: resource.MustParse("128m"),
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "ssh",
					ReadOnly:  true,
					MountPath: "/root/.ssh",
				},
				{
					Name:      "rootcomposer",
					MountPath: "/root/.composer",
				},
				{
					Name:      GitCloneVolume,
					MountPath: "/data",
				},
				{
					Name:      "tmp",
					MountPath: "/tmp",
				},
			},
		},
	}, have.Spec.Containers)

	assert.Equal(t, []corev1.Volume{
		{
			Name: "ssh",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "test-ssh",
					DefaultMode: &perm,
				},
			},
		},
		{
			Name: GitCloneVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumDefault,
				},
			},
		},
		{
			Name: "rootcomposer",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "rootcomposer",
				},
			},
		},
		{
			Name: "tmp",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumMemory,
				},
			},
		},
	}, have.Spec.Volumes)

	assert.Equal(t, []corev1.LocalObjectReference{
		{
			Name: "test-dockercfg",
		},
	}, have.Spec.ImagePullSecrets)
}
