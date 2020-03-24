package env

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/previousnext/m8s/pb"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TypeMySQL is used of idenifying a MySQL container for readiness rules.
const TypeMySQL = "mysql"

// PodInput provides the Pod function with information to produce a Kubernetes Pod.
type PodInput struct {
	Namespace       string
	Name            string
	Annotations     []*pb.Annotation
	Repository      string
	Revision        string
	Retention       string
	Services        []*pb.ComposeService
	Caches          []PodInputCache
	ImagePullSecret string
	Init            []*pb.Init
	Domain string
}

// PodInputCache is used for passing in cache configuration to generate a pod.
type PodInputCache struct {
	Name string
	Path string
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
				"author": "m8s",
				"rig.io/domain": input.Domain,
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
			},
		},
	}

	// These container will clone the code into an "emptyDir" volume.
	cloneContainers, cloneVolume := GitCloneInitContainers(input.Repository, input.Revision)
	pod.Spec.InitContainers = cloneContainers
	pod.Spec.Volumes = append(pod.Spec.Volumes, cloneVolume)

	if input.ImagePullSecret != "" {
		pod.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: input.ImagePullSecret,
			},
		}
	}

	for _, cache := range input.Caches {
		pod.Spec.Volumes = append(pod.Spec.Volumes, v1.Volume{
			Name: cache.Name,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: cache.Name,
				},
			},
		})
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

	for _, init := range input.Init {
		container := v1.Container{
			Image:           init.Image,
			WorkingDir:      GitClonePath,
			ImagePullPolicy: v1.PullAlways,
			VolumeMounts: []v1.VolumeMount{
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
		}

		resources, err := podResources(init.Reservations, init.Limits)
		if err != nil {
			return pod, err
		}
		container.Resources = resources

		for _, cache := range input.Caches {
			container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
				Name:      cache.Name,
				MountPath: cache.Path,
			})
		}

		for stepID, stepCommand := range init.Steps {
			container.Name = fmt.Sprintf("%s-step%d", init.Name, stepID+1)
			container.Command = strings.Split(stepCommand, " ")

			pod.Spec.InitContainers = append(pod.Spec.InitContainers, container)
		}
	}

	hostAliases := make(map[string][]string)

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
			},
		}

		if service.Type == TypeMySQL {
			container.ReadinessProbe = &v1.Probe{
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: []string{
							"mysqladmin", "ping", "-h", "127.0.0.1",
						},
					},
				},
				InitialDelaySeconds: 15,
				TimeoutSeconds:      30,
				PeriodSeconds:       15,
				SuccessThreshold:    1,
				FailureThreshold:    10,
			}
		}

		if len(service.Entrypoint) > 0 {
			container.Command = service.Entrypoint
		}

		resources, err := podResources(service.Reservations, service.Limits)
		if err != nil {
			return pod, err
		}

		mounts, volumes, err := podVolumes(service.Volumes, service.Tmpfs, input.Caches)
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


		if len(service.Extrahosts) > 0 {
			for _, value := range service.Extrahosts {
				parts := strings.Split(value, ":")
				hostAliases[parts[1]] = append(hostAliases[parts[1]], parts[0])
			}
		}
	}

	for ip, hostnames := range hostAliases {
		pod.Spec.HostAliases = append(pod.Spec.HostAliases, v1.HostAlias{
			IP:        ip,
			Hostnames: unique(hostnames),
		})
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
func podVolumes(serviceVolumes []string, tmps []string, caches []PodInputCache) ([]v1.VolumeMount, []v1.Volume, error) {
	var (
		mounts  []v1.VolumeMount
		volumes []v1.Volume
	)

	for _, cache := range caches {
		mounts = append(mounts, v1.VolumeMount{
			Name:      cache.Name,
			MountPath: cache.Path,
		})
	}

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
				Name:      GitCloneVolume,
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

// Helper function to return a slice of unique values.
func unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
