package openshift

import (
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	routev1client *routev1.RouteV1Client
	client *kubernetes.Clientset
	config *rest.Config
}