package pod

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/previousnext/compose"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/m8s/utils"
	"github.com/previousnext/m8s/utils/k8s/pod/sidecar"
)

// GenerateParams for generating an Pod object.
type GenerateParams struct {
	Namespace       string
	Name            string
	Annotations     map[string]string
	Repository      string
	Revision        string
	Services        map[string]compose.Service
	Caches          []string
	SecretDockerCfg string
	SecretSSH       string
	Sidecar         sidecar.GenerateParams
}

// Generate will generate a Pod object.
func Generate(params GenerateParams) (*corev1.Pod, error) {
	// Permissions value used by SSH id_rsa key.
	// https://kubernetes.io/docs/user-guide/secrets/
	perm := int32(256)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: params.Namespace,
			Name:      params.Name,
			// This allows us to Link our Service to this Pod.
			Labels: map[string]string{
				"env": params.Name,
			},
			Annotations: params.Annotations,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "ssh",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  params.SecretSSH,
							DefaultMode: &perm,
						},
					},
				},
			},
		},
	}

	// These container will clone the code into an "emptyDir" volume.
	cloneContainers, cloneVolume := GitCloneInitContainers(params.Repository, params.Revision)
	pod.Spec.InitContainers = cloneContainers
	pod.Spec.Volumes = append(pod.Spec.Volumes, cloneVolume)

	// This container is used for routing.
	s, err := sidecar.Generate(params.Sidecar)
	if err != nil {
		return pod, errors.Wrap(err, "failed to generate sidecar")
	}
	pod.Spec.Containers = append(pod.Spec.Containers, s)

	if params.SecretDockerCfg != "" {
		pod.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: params.SecretDockerCfg,
			},
		}
	}

	for _, cache := range params.Caches {
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: utils.Machine(cache),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: utils.Machine(cache),
				},
			},
		})
	}

	for name, service := range params.Services {
		container := corev1.Container{
			Name:            name,
			Image:           service.Image,
			ImagePullPolicy: corev1.PullAlways,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "ssh",
					ReadOnly:  true,
					MountPath: "/root/.ssh",
				},
			},
		}

		if len(service.Entrypoint) > 0 {
			container.Command = service.Entrypoint
		}

		resources, err := podResources(service.Deploy.Resources.Reservations, service.Deploy.Resources.Limits)
		if err != nil {
			return pod, err
		}

		mounts, volumes, err := generatePodVolumes(service.Volumes, service.Tmpfs, params.Caches)
		if err != nil {
			return pod, err
		}

		ports, err := generatePodPorts(service.Ports)
		if err != nil {
			return pod, err
		}

		envs, err := generatePodEnvs(service.Environment)
		if err != nil {
			return pod, err
		}

		securityContext, err := generatePodSecurity(service.CapAdd)
		if err != nil {
			return pod, err
		}

		if len(service.CapAdd) > 0 {
			container.SecurityContext = securityContext
		}

		container.Resources = resources
		container.VolumeMounts = append(container.VolumeMounts, mounts...)
		container.Ports = append(container.Ports, ports...)
		container.Env = append(container.Env, envs...)

		// Add volumes and containers to the pod definition.
		pod.Spec.Volumes = append(pod.Spec.Volumes, volumes...)
		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	return pod, nil
}

// Helper function to extract resource limits from a service definition.
func podResources(reservations, limits compose.ServiceDeployResource) (corev1.ResourceRequirements, error) {
	resources := corev1.ResourceRequirements{
		Limits:   make(map[corev1.ResourceName]resource.Quantity),
		Requests: make(map[corev1.ResourceName]resource.Quantity),
	}

	if reservations.CPUs != "" {
		quantity, err := resource.ParseQuantity(reservations.CPUs)
		if err != nil {
			return resources, fmt.Errorf("failed to parse cpu reservation: %s", err)
		}

		resources.Requests[corev1.ResourceCPU] = quantity
	}

	if reservations.Memory != "" {
		quantity, err := resource.ParseQuantity(reservations.Memory)
		if err != nil {
			return resources, fmt.Errorf("failed to parse memory reservation: %s", err)
		}

		resources.Requests[corev1.ResourceMemory] = quantity
	}

	if limits.CPUs != "" {
		quantity, err := resource.ParseQuantity(limits.CPUs)
		if err != nil {
			return resources, fmt.Errorf("failed to parse cpu limits: %s", err)
		}

		resources.Limits[corev1.ResourceCPU] = quantity
	}

	if limits.Memory != "" {
		quantity, err := resource.ParseQuantity(limits.Memory)
		if err != nil {
			return resources, fmt.Errorf("failed to parse memory limits: %s", err)
		}

		resources.Limits[corev1.ResourceMemory] = quantity
	}

	return resources, nil
}

// Helper function to extract volumes from a service definition.
func generatePodVolumes(serviceVolumes []string, tmps []string, caches []string) ([]corev1.VolumeMount, []corev1.Volume, error) {
	var (
		mounts  []corev1.VolumeMount
		volumes []corev1.Volume
	)

	for _, cache := range caches {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      utils.Machine(cache),
			MountPath: cache,
		})
	}

	// Adds the Docker Compose volumes to our Pod object.
	exists, path := FindCodePath(serviceVolumes)
	if exists {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      GitCloneVolume,
			MountPath: path,
		})
	}

	for _, tmp := range tmps {
		name := utils.Machine(tmp)

		volumes = append(volumes, corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumMemory,
				},
			},
		})

		mounts = append(mounts, corev1.VolumeMount{
			Name:      name,
			MountPath: tmp,
		})
	}

	return mounts, volumes, nil
}

// Helper function to extract ports from a service definition.
func generatePodPorts(list []string) ([]corev1.ContainerPort, error) {
	var ports []corev1.ContainerPort

	for _, item := range list {
		sl := strings.Split(item, ":")

		// Ensure we have an environment variable in the format "FOO=bar".
		if len(sl) < 1 {
			continue
		}

		val, err := strconv.ParseInt(sl[0], 10, 32)
		if err != nil {
			continue
		}

		ports = append(ports, corev1.ContainerPort{
			ContainerPort: int32(val),
		})
	}

	return ports, nil
}

// Helper function to extract environment variables from a service definition.
func generatePodEnvs(list []string) ([]corev1.EnvVar, error) {
	var envs []corev1.EnvVar

	for _, item := range list {
		sl := strings.Split(item, "=")

		// Ensure we have an environment variable in the format "FOO=bar".
		if len(sl) != 2 {
			continue
		}

		envs = append(envs, corev1.EnvVar{
			Name:  sl[0],
			Value: sl[1],
		})
	}

	return envs, nil
}

// Helper function to extract a security context for a container.
func generatePodSecurity(adds []string) (*corev1.SecurityContext, error) {
	caps := &corev1.Capabilities{}

	for _, add := range adds {
		caps.Add = append(caps.Add, corev1.Capability(add))
	}

	return &corev1.SecurityContext{
		Capabilities: caps,
	}, nil
}
