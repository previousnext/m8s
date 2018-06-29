package main

import (
	"log"
	"time"

	"github.com/alecthomas/kingpin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

var (
	namespace = kingpin.Flag("namespace", "Namespace to scan for environments").Default(corev1.NamespaceAll).OverrideDefaultFromEnvar("NAMESPACE").String()
)

func main() {
	kingpin.Parse()

	limiter := time.Tick(time.Minute)

	for {
		<-limiter

		now := time.Now().Unix()

		// Creates the in-cluster config.
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println("Failed to get cluster configuration:", err)
			continue
		}

		// Creates the clientset for querying APIs.
		k8s, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Println("Failed to get client:", err)
			continue
		}

		log.Println("Killing Ingresses")

		err = killIngresses(k8s, *namespace, now)
		if err != nil {
			log.Println("Failed to kill Ingress:", err)
		}

		log.Println("Killing Services")

		err = killServices(k8s, *namespace, now)
		if err != nil {
			log.Println("Failed to kill Services:", err)
		}

		log.Println("Killing Pods")

		err = killPods(k8s, *namespace, now)
		if err != nil {
			log.Println("Failed to kill Pods:", err)
		}

		log.Println("Killing ReplicaSets")

		err = killReplicaSets(k8s, *namespace, now)
		if err != nil {
			log.Println("Failed to kill ReplicaSets:", err)
		}

		log.Println("Killing Deployments")

		err = killDeployments(k8s, *namespace, now)
		if err != nil {
			log.Println("Failed to kill Deployments:", err)
		}
	}
}

// Helper function to deleting old Ingress objects.
func killIngresses(clientset *kubernetes.Clientset, namespace string, now int64) error {
	ings, err := clientset.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ing := range ings.Items {
		bd, err := getBlackDeath(ing.ObjectMeta.Annotations)
		if err != nil {
			log.Println("Skipping Ingress:", ing.Name, err)
			continue
		}

		if bd < now {
			log.Println("Deleting Ingress:", ing.Name)

			err := clientset.ExtensionsV1beta1().Ingresses(ing.Namespace).Delete(ing.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to deleting old Ingress objects.
func killServices(clientset *kubernetes.Clientset, namespace string, now int64) error {
	svcs, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, svc := range svcs.Items {
		bd, err := getBlackDeath(svc.ObjectMeta.Annotations)
		if err != nil {
			log.Println("Skipping Service:", svc.Name, err)
			continue
		}

		if bd < now {
			log.Println("Deleting Service:", svc.Name)

			err := clientset.CoreV1().Services(svc.Namespace).Delete(svc.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to deleting old Ingress objects.
func killPods(clientset *kubernetes.Clientset, namespace string, now int64) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		bd, err := getBlackDeath(pod.ObjectMeta.Annotations)
		if err != nil {
			log.Println("Skipping Pod:", pod.Name, err)
			continue
		}

		if bd < now {
			log.Println("Deleting Pod:", pod.Name)

			err := clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to deleting old ReplicaSet objects.
func killReplicaSets(clientset *kubernetes.Clientset, namespace string, now int64) error {
	sets, err := clientset.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, set := range sets.Items {
		bd, err := getBlackDeath(set.ObjectMeta.Annotations)
		if err != nil {
			log.Println("Skipping ReplicaSet:", set.Name, err)
			continue
		}

		if bd < now {
			log.Println("Deleting ReplicaSet:", set.Name)

			err := clientset.ExtensionsV1beta1().ReplicaSets(set.Namespace).Delete(set.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to deleting old Deployment objects.
func killDeployments(clientset *kubernetes.Clientset, namespace string, now int64) error {
	dplys, err := clientset.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, dply := range dplys.Items {
		bd, err := getBlackDeath(dply.ObjectMeta.Annotations)
		if err != nil {
			log.Println("Skipping Deployment:", dply.Name, err)
			continue
		}

		if bd < now {
			log.Println("Deleting Deployment:", dply.Name)

			err := clientset.ExtensionsV1beta1().Deployments(dply.Namespace).Delete(dply.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
