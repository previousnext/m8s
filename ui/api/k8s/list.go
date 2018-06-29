package k8s

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/previousnext/m8s/ui/api/types"
	apiutils "github.com/previousnext/m8s/ui/api/utils"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// List returns a list of environments.
func (s Server) List(w http.ResponseWriter, r *http.Request) {
	config, err := clientcmd.BuildConfigFromFlags(s.Master, s.Config)
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	ingresses, err := client.ExtensionsV1beta1().Ingresses(s.Namespace).List(metav1.ListOptions{})
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	pods, err := client.CoreV1().Pods(s.Namespace).List(metav1.ListOptions{})
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	var list []types.Environment

	for _, ingress := range ingresses.Items {
		for _, rule := range ingress.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Path == "/" {
					// Check if the pod which the same name as this service exists.
					pod, err := findPod(path.Backend.ServiceName, pods.Items)
					if err != nil {
						continue
					}

					item := types.Environment{
						Name:   pod.Name,
						Domain: rule.Host,
					}

					for _, container := range pod.Spec.Containers {
						item.Containers = append(item.Containers, container.Name)
					}

					list = append(list, item)
				}
			}
		}
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		apiutils.Fatal(w, err)
	}
}

func findPod(needle string, haystack []v1.Pod) (v1.Pod, error) {
	for _, pod := range haystack {
		if pod.Name == needle {
			return pod, nil
		}
	}

	return v1.Pod{}, fmt.Errorf("pod with name %s does not exist", needle)
}
