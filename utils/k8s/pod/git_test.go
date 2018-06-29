package pod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestGitCloneInitContainers(t *testing.T) {
	containers, volume := GitCloneInitContainers("https://github.com/example/repo.git", "acb123")

	assert.Equal(t, []corev1.Container{
		{
			Name:            "git-init",
			Image:           GitCloneImage,
			ImagePullPolicy: corev1.PullAlways,
			WorkingDir:      GitClonePath,
			Command: []string{
				"git",
				"init",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      GitCloneVolume,
					MountPath: GitClonePath,
				},
			},
		},
		{
			Name:            "git-remote",
			Image:           GitCloneImage,
			ImagePullPolicy: corev1.PullAlways,
			WorkingDir:      GitClonePath,
			Command: []string{
				"git",
				"remote",
				"add",
				"origin",
				"https://github.com/example/repo.git",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      GitCloneVolume,
					MountPath: GitClonePath,
				},
			},
		},
		{
			Name:            "git-fetch",
			Image:           GitCloneImage,
			ImagePullPolicy: corev1.PullAlways,
			WorkingDir:      GitClonePath,
			Command: []string{
				"git",
				"fetch",
				"origin",
				"acb123",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      GitCloneVolume,
					MountPath: GitClonePath,
				},
			},
		},
		{
			Name:            "git-reset",
			Image:           GitCloneImage,
			ImagePullPolicy: corev1.PullAlways,
			WorkingDir:      GitClonePath,
			Command: []string{
				"git",
				"reset",
				"--hard",
				"FETCH_HEAD",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      GitCloneVolume,
					MountPath: GitClonePath,
				},
			},
		},
	}, containers)

	assert.Equal(t, corev1.Volume{
		Name: GitCloneVolume,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: corev1.StorageMediumDefault,
			},
		},
	}, volume)
}
