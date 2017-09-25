package utils

import (
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

// PodCreate is used for creating a Pod object.
func PodCreate(client *client.Clientset, pod *v1.Pod) (*v1.Pod, error) {
	_, err := client.Pods(pod.ObjectMeta.Namespace).Create(pod)
	if errors.IsAlreadyExists(err) {
		// This will tell Kubernetes that we want this pod to be deleted immediately.
		now := int64(0)

		// Delete the Pod.
		err = client.Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{
			GracePeriodSeconds: &now,
		})
		if err != nil {
			return pod, err
		}

		// Create the new pod.
		_, err = client.Pods(pod.ObjectMeta.Namespace).Create(pod)
		if err != nil {
			return pod, err
		}
	} else if err != nil {
		return pod, err
	}

	// Wait for the pod to become available.
	limiter := time.Tick(time.Second / 10)

	for {
		pod, err = client.Pods(pod.ObjectMeta.Namespace).Get(pod.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return pod, err
		}

		if pod.Status.Phase == v1.PodRunning {
			break
		}

		<-limiter
	}

	return pod, err
}
