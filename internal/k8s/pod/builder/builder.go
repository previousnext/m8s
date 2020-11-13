package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// AnnotationDomain which is used to routing traffic to this instance.
	AnnotationDomain = "m8s.io/domain"
)

type GenerateParams struct {
	Name string
	Namespace string
	Domain string
	// ExtraAnnotations which are applied to the Pod eg. CircleCI/Github metadata from a build pipeline.
	ExtraAnnotations map[string]string
	ExtraEnvironment map[string]string
}

// Generate a Pod.
func Generate(params GenerateParams) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Name,
			Namespace: params.Namespace,
			Annotations: map[string]string{
				// @todo, K8s standard annotations.
				AnnotationDomain: params.Domain,
			},
		},
	}

	for name, value := range params.ExtraAnnotations {
		pod.ObjectMeta.Annotations[name] = value
	}

	if err := applyVolumes(pod); err != nil {
		return nil, err
	}

	if err := applyInitContainers(pod); err != nil {
		return nil, err
	}

	if err := applyContainers(pod); err != nil {
		return nil, err
	}

	if err := applyHostAliases(pod); err != nil {
		return nil, err
	}

	return pod, nil
}

func applyVolumes(pod *corev1.Pod) error {
	return nil
}

func applyInitContainers(pod *corev1.Pod) error {
	// Apply Git clone
	// Apply build steps
	return nil
}

func applyContainers(pod *corev1.Pod) error {
	return nil
}

func applyHostAliases(pod *corev1.Pod) error {
	return nil
}
