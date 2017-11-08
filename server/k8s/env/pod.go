package env

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/previousnext/m8s/pb"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// PodInput provides the Pod function with information to produce a Kubernetes Pod.
type PodInput struct {
	Namespace   string
	Name        string
	Annotations []*pb.Annotation
	Repository  string
	Revision    string
	Retention   string
	Services    []*pb.ComposeService
	Prometheus  int32
}

// Pod converts a Docker Compose file into a Kubernetes Deployment object.
func Pod(input PodInput) (*v1.Pod, error) {
	// Permissions value used by SSH id_rsa key.
	// https://kubernetes.io/docs/user-guide/secrets/
	perm := int32(256)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: input.Namespace,
			Name:      input.Name,
			// This allows us to Link our Service to this Pod.
			Labels: map[string]string{
				"env": input.Name,
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   fmt.Sprintf("%d", input.Prometheus),
				"author":               "m8s",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "apache-exporter",
					Image: "previousnext/apache-exporter:latest",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: input.Prometheus,
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
							Repository: input.Repository,
							Revision:   input.Revision,
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

	for _, annotation := range input.Annotations {
		pod.ObjectMeta.Annotations[annotation.Name] = annotation.Value
	}

	if input.Retention != "" {
		unix, err := retentionToUnix(time.Now(), input.Retention)
		if err != nil {
			return pod, err
		}

		pod.ObjectMeta.Annotations["black-death.skpr.io"] = unix
	}

	for _, service := range input.Services {
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

		resources, err := podResources(service.Reservations, service.Limits)
		if err != nil {
			return pod, err
		}

		mounts, volumes, err := podVolumes(service.Volumes, service.Tmpfs)
		if err != nil {
			return pod, err
		}

		ports, err := podPorts(service.Ports)
		if err != nil {
			return pod, err
		}

		envs, err := podEnvs(service.Environment)
		if err != nil {
			return pod, err
		}

		securityContext, err := podSecurity(service.Capabilities)
		if err != nil {
			return pod, err
		}

		container.Resources = resources
		container.VolumeMounts = append(container.VolumeMounts, mounts...)
		container.Ports = append(container.Ports, ports...)
		container.Env = append(container.Env, envs...)

		if len(service.Capabilities) > 0 {
			container.SecurityContext = securityContext
		}

		// Add volumes and containers to the pod definition.
		pod.Spec.Volumes = append(pod.Spec.Volumes, volumes...)
		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	return pod, nil
}

// Helper function to extract resource limits from a service definition.
func podResources(reservations, limits *pb.Resource) (v1.ResourceRequirements, error) {
	resources := v1.ResourceRequirements{
		Limits:   make(map[v1.ResourceName]resource.Quantity),
		Requests: make(map[v1.ResourceName]resource.Quantity),
	}

	if reservations != nil && reservations.CPU != "" {
		quantity, err := resource.ParseQuantity(reservations.CPU)
		if err != nil {
			return resources, fmt.Errorf("failed to parse cpu reservation: %s", err)
		}

		resources.Requests[v1.ResourceCPU] = quantity
	}

	if reservations != nil && reservations.Memory != "" {
		quantity, err := resource.ParseQuantity(reservations.Memory)
		if err != nil {
			return resources, fmt.Errorf("failed to parse memory reservation: %s", err)
		}

		resources.Requests[v1.ResourceMemory] = quantity
	}

	if limits != nil && limits.CPU != "" {
		quantity, err := resource.ParseQuantity(limits.CPU)
		if err != nil {
			return resources, fmt.Errorf("failed to parse cpu limits: %s", err)
		}

		resources.Limits[v1.ResourceCPU] = quantity
	}

	if limits != nil && limits.Memory != "" {
		quantity, err := resource.ParseQuantity(limits.Memory)
		if err != nil {
			return resources, fmt.Errorf("failed to parse memory limits: %s", err)
		}

		resources.Limits[v1.ResourceMemory] = quantity
	}

	return resources, nil
}

// Helper function to extract volumes from a service definition.
func podVolumes(serviceVolumes []string, tmps []string) ([]v1.VolumeMount, []v1.Volume, error) {
	var (
		mounts  []v1.VolumeMount
		volumes []v1.Volume
	)

	// Adds the Docker Compose volumes to our Pod object.
	for _, serviceVolume := range serviceVolumes {
		sl := strings.Split(serviceVolume, ":")

		// Ensure we have an volume in the format "/home/nick:/data".
		if len(sl) < 2 {
			continue
		}

		// Mount the code where the user has provided "." as the "source".
		// Anything else the user has provided cannot be supported.
		// @todo, Handle other mounts.
		if sl[0] == "." {
			mounts = append(mounts, v1.VolumeMount{
				Name:      "code",
				MountPath: sl[1],
			})
		}
	}

	for _, tmp := range tmps {
		name := machine(tmp)

		volumes = append(volumes, v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium: v1.StorageMediumMemory,
				},
			},
		})

		mounts = append(mounts, v1.VolumeMount{
			Name:      name,
			MountPath: tmp,
		})
	}

	return mounts, volumes, nil
}

// Helper function to extract ports from a service definition.
func podPorts(list []string) ([]v1.ContainerPort, error) {
	var ports []v1.ContainerPort

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

		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(val),
		})
	}

	return ports, nil
}

// Helper function to extract environment variables from a service definition.
func podEnvs(list []string) ([]v1.EnvVar, error) {
	var envs []v1.EnvVar

	for _, item := range list {
		sl := strings.Split(item, "=")

		// Ensure we have an environment variable in the format "FOO=bar".
		if len(sl) != 2 {
			continue
		}

		envs = append(envs, v1.EnvVar{
			Name:  sl[0],
			Value: sl[1],
		})
	}

	return envs, nil
}

// Helper function to extract a securit context for a container.
func podSecurity(adds []string) (*v1.SecurityContext, error) {
	caps := &v1.Capabilities{}

	for _, add := range adds {
		caps.Add = append(caps.Add, v1.Capability(add))
	}

	return &v1.SecurityContext{
		Capabilities: caps,
	}, nil
}
