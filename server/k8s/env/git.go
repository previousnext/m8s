package env

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// GitCloneImage is the image used for cloning the repository.
	GitCloneImage = "alpine/git:latest"
	// GitCloneVolume is the name of the volume which code is cloned into.
	GitCloneVolume = "code"
	// GitClonePath is the path which the init containers mount the GitCloneVolume.
	GitClonePath = "/checkout"
)

// GitCloneInitContainers returns a list of init containers and volume which they clone a repository into.
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
					{
						Name:      SecretSSH,
						ReadOnly:  true,
						MountPath: "/root/.ssh",
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
					{
						Name:      SecretSSH,
						ReadOnly:  true,
						MountPath: "/root/.ssh",
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
					{
						Name:      SecretSSH,
						ReadOnly:  true,
						MountPath: "/root/.ssh",
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
					{
						Name:      SecretSSH,
						ReadOnly:  true,
						MountPath: "/root/.ssh",
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
