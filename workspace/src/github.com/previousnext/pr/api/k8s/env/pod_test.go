package env

import (
	"testing"

	pb "github.com/previousnext/pr/pb"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

func TestPod(t *testing.T) {
	perm := int32(256)

	want := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
			Labels: map[string]string{
				"env": "pr1",
			},
			Annotations: map[string]string{
				"skipper.io/black-death": "123456789",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "app",
					Image:           "foo/bar",
					ImagePullPolicy: v1.PullAlways,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
						{
							Name:      "code",
							MountPath: "/data",
						},
					},
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
					Env: []v1.EnvVar{
						{
							Name:  "FOO",
							Value: "bar",
						},
					},
				},
				{
					Name:            "mysql",
					Image:           "mariadb",
					ImagePullPolicy: v1.PullAlways,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
					},
					Env: []v1.EnvVar{
						{
							Name:  "MYSQL_ROOT_PASSWORD",
							Value: "root",
						},
						{
							Name:  "MYSQL_DATABASE",
							Value: "local",
						},
						{
							Name:  "MYSQL_USER",
							Value: "drupal",
						},
						{
							Name:  "MYSQL_PASSWORD",
							Value: "drupal",
						},
					},
				},
				{
					Name:            "solr",
					Image:           "previousnext/solr:5.x",
					ImagePullPolicy: v1.PullAlways,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "ssh",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName:  SecretSSH,
							DefaultMode: &perm,
						},
					},
				},
				{
					Name: "code",
					VolumeSource: v1.VolumeSource{
						GitRepo: &v1.GitRepoVolumeSource{
							Repository: "git@github.com:foo/bar.git",
							Revision:   "123456789",
							Directory:  ".",
						},
					},
				},
			},
			ImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: SecretDockerCfg,
				},
			},
		},
	}

	have, err := Pod(123456789, "test", "pr1", "git@github.com:foo/bar.git", "123456789", []*pb.ComposeService{
		{
			Name:  "app",
			Image: "foo/bar",
			Volumes: []string{
				".:/data",
			},
			Ports: []string{
				"80:80",
			},
			Environment: []string{
				"FOO=bar",
			},
		},
		{
			Name:  "mysql",
			Image: "mariadb",
			Environment: []string{
				"MYSQL_ROOT_PASSWORD=root",
				"MYSQL_DATABASE=local",
				"MYSQL_USER=drupal",
				"MYSQL_PASSWORD=drupal",
			},
		},
		{
			Name:  "solr",
			Image: "previousnext/solr:5.x",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, want.ObjectMeta, have.ObjectMeta)
	assert.Equal(t, want.Spec.Containers, have.Spec.Containers)
	assert.Equal(t, want.Spec.Volumes, have.Spec.Volumes)
	assert.Equal(t, want.Spec.ImagePullSecrets, have.Spec.ImagePullSecrets)
}
