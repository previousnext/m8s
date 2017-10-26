package env

import (
	"testing"

	"github.com/previousnext/m8s/cmd/metadata"
	pb "github.com/previousnext/m8s/pb"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func TestPod(t *testing.T) {
	var (
		perm = int32(256)
		prom = int32(9117)
	)

	want := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
			Labels: map[string]string{
				"env": "pr1",
			},
			Annotations: map[string]string{
				"author":                              "m8s",
				"prometheus.io/scrape":                "true",
				"prometheus.io/port":                  "9117",
				metadata.AnnotationBitbucketRepoOwner: "nick",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "apache-exporter",
					Image: "previousnext/apache-exporter:latest",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 9117,
						},
					},
				},
				{
					Name:            "app",
					Image:           "foo/bar",
					ImagePullPolicy: v1.PullAlways,
					Resources: v1.ResourceRequirements{
						Limits:   v1.ResourceList{},
						Requests: v1.ResourceList{},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
						{
							Name:      CacheComposer,
							MountPath: "/root/.composer",
						},
						{
							Name:      CacheYarn,
							MountPath: "/usr/local/share/.cache/yarn",
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
					Resources: v1.ResourceRequirements{
						Limits:   v1.ResourceList{},
						Requests: v1.ResourceList{},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
						{
							Name:      CacheComposer,
							MountPath: "/root/.composer",
						},
						{
							Name:      CacheYarn,
							MountPath: "/usr/local/share/.cache/yarn",
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
					Resources: v1.ResourceRequirements{
						Limits:   v1.ResourceList{},
						Requests: v1.ResourceList{},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      SecretSSH,
							ReadOnly:  true,
							MountPath: "/root/.ssh",
						},
						{
							Name:      CacheComposer,
							MountPath: "/root/.composer",
						},
						{
							Name:      CacheYarn,
							MountPath: "/usr/local/share/.cache/yarn",
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
				{
					Name: CacheComposer,
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: CacheComposer,
						},
					},
				},
				{
					Name: CacheYarn,
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: CacheYarn,
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

	annotations, err := metadata.Annotations([]string{"BITBUCKET_REPO_OWNER=nick"})
	assert.Nil(t, err)

	have, err := Pod("test", "pr1", annotations, "git@github.com:foo/bar.git", "123456789", "", []*pb.ComposeService{
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
	}, prom)
	assert.Nil(t, err)
	assert.Equal(t, want.ObjectMeta, have.ObjectMeta)
	assert.Equal(t, want.Spec.Containers[0], have.Spec.Containers[0])
	assert.Equal(t, want.Spec.Containers[1], have.Spec.Containers[1])
	assert.Equal(t, want.Spec.Containers[2], have.Spec.Containers[2])
	assert.Equal(t, want.Spec.Volumes, have.Spec.Volumes)
	assert.Equal(t, want.Spec.ImagePullSecrets, have.Spec.ImagePullSecrets)
}
