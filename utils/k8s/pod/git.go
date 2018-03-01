package pod

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	GitCloneImage  = "alpine/git:latest"
	GitCloneVolume = "code"
	GitClonePath   = "/checkout"
)

func GitCloneInitContainers(repository, revision string) ([]corev1.Container, corev1.Volume) {
	return []corev1.Container{
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
					repository,
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
					revision,
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
		}, corev1.Volume{
			Name: GitCloneVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumDefault,
				},
			},
		}
}
