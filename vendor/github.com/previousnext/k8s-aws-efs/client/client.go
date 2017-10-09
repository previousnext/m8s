package client

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	API      = "efs"
	Group    = "skpr.io"
	Resource = "efses"
)

type Client struct {
	rc *rest.RESTClient
}

// Create a client for Efs interactions.
func NewClient(config *rest.Config) (Client, error) {
	var c Client

	rc, err := rest.RESTClientFor(solrConfig(config))
	if err != nil {
		return c, err
	}

	// Store the rest client for future queries.
	// We will use this client for Get() and List() operations.
	c.rc = rc

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return c, err
	}

	// Create a K8s Third Party Resource object.
	// This allows us to start creating Efs filesystems under a bespoke, K8s backed, API.
	_, err = clientset.Extensions().ThirdPartyResources().Create(&v1beta1.ThirdPartyResource{
		ObjectMeta: metav1.ObjectMeta{
			Name: API + "." + Group,
		},
		Versions: []v1beta1.APIVersion{
			{Name: "v1"},
		},
		Description: "A ThirdPartyResource for provisioning AWS EFS resources",
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return c, err
	}

	return c, nil
}

// Returns a single Sore core for a namespace.
func (c *Client) Get(namespace, name string) (Efs, error) {
	var s Efs
	err := c.rc.Get().Resource(Resource).Namespace(namespace).Name(name).Do().Into(&s)
	return s, err
}

// Sets the entire EFS object (Spec + Status).
func (c *Client) Put(efs Efs) error {
	return c.rc.Put().Resource(Resource).Namespace(efs.Metadata.Namespace).Name(efs.Metadata.Name).Body(&efs).Do().Error()
}

// Sets the entire EFS object (Spec + Status).
func (c *Client) Post(efs Efs) error {
	return c.rc.Post().Resource(Resource).Namespace(efs.Metadata.Namespace).Body(&efs).Do().Error()
}

// Returns a list of Solr cores from all namespaces.
func (c *Client) List(namespace string) (EfsList, error) {
	s := EfsList{}
	err := c.rc.Get().Resource(Resource).Namespace(namespace).Do().Into(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// Returns a list of Solr cores from all namespaces.
func (c *Client) ListAll() (EfsList, error) {
	s := EfsList{}
	err := c.rc.Get().Resource(Resource).Do().Into(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
