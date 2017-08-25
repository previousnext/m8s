package env

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/previousnext/m8s/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Pod converts a Docker Compose file into a Kubernetes Deployment object.
func Pod(timeout int64, namespace, name, repository, revision string, services []*pb.ComposeService) (*v1.Pod, error) {
	// Permissions value used by SSH id_rsa key.
	// https://kubernetes.io/docs/user-guide/secrets/
	perm := int32(256)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			// This allows us to Link our Service to this Pod.
			Labels: map[string]string{
				"env": name,
			},
			Annotations: map[string]string{
				"skipper.io/black-death": fmt.Sprintf("%v", timeout),
			},
		},
		Spec: v1.PodSpec{
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
							Repository: repository,
							Revision:   revision,
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

	for _, service := range services {
		container := v1.Container{
			Name:            service.Name,
			Image:           service.Image,
			ImagePullPolicy: v1.PullAlways,
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
		}

		// Adds the Docker Compose volumes to our Pod object.
		for _, volume := range service.Volumes {
			sl := strings.Split(volume, ":")

			// Ensure we have an volume in the format "/home/nick:/data".
			if len(sl) < 2 {
				continue
			}

			// Mount the code where the user has provided "." as the "source".
			// Anything else the user has provided cannot be supported.
			// @todo, Handle other mounts.
			if sl[0] == "." {
				container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
					Name:      "code",
					MountPath: sl[1],
				})
			}
		}

		// Adds the Docker Compose ports to our Pod object.
		for _, port := range service.Ports {
			sl := strings.Split(port, ":")

			// Ensure we have an environment variable in the format "FOO=bar".
			if len(sl) < 1 {
				continue
			}

			val, err := strconv.ParseInt(sl[0], 10, 32)
			if err != nil {
				continue
			}

			container.Ports = append(container.Ports, v1.ContainerPort{
				ContainerPort: int32(val),
			})
		}

		// Adds the Docker Compose environment variables to our Pod object.
		for _, env := range service.Environment {
			sl := strings.Split(env, "=")

			// Ensure we have an environment variable in the format "FOO=bar".
			if len(sl) != 2 {
				continue
			}

			container.Env = append(container.Env, v1.EnvVar{
				Name:  sl[0],
				Value: sl[1],
			})
		}

		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	return pod, nil
}
