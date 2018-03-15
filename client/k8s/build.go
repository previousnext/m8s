package k8s

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/m8s/client/types"
	"github.com/previousnext/m8s/utils"
	m8singress "github.com/previousnext/m8s/utils/k8s/ingress"
	m8spod "github.com/previousnext/m8s/utils/k8s/pod"
	"github.com/previousnext/m8s/utils/k8s/pod/sidecar"
	m8sclaim "github.com/previousnext/m8s/utils/k8s/pvc"
	m8sservice "github.com/previousnext/m8s/utils/k8s/service"
	skpringress "github.com/previousnext/skpr/utils/k8s/ingress"
	skprpod "github.com/previousnext/skpr/utils/k8s/pods"
	skprclaim "github.com/previousnext/skpr/utils/k8s/pvc"
	skprservice "github.com/previousnext/skpr/utils/k8s/service"
)

// Build the environment.
func (c Client) Build(w io.Writer, params types.BuildParams) error {
	for _, path := range params.Config.Cache.Paths {
		fmt.Fprintf(w, "Creating: PersistentVolumeClaim: %s\n", path)

		err := createClaim(c.client, params, path)
		if err != nil {
			return errors.Wrap(err, "failed to create PersistentVolumeClaim")
		}
	}

	fmt.Fprintln(w, "Creating: Service")

	err := createService(c.client, params)
	if err != nil {
		return errors.Wrap(err, "failed to create Service")
	}

	fmt.Fprintln(w, "Creating: Ingress")

	err = createIngress(c.client, params)
	if err != nil {
		return errors.Wrap(err, "failed to create Ingress")
	}

	fmt.Fprintln(w, "Creating: Pod")

	err = createPod(c.client, params)
	if err != nil {
		return errors.Wrap(err, "failed to create Pod")
	}

	return nil
}

// Helper to create a PersistentVolumeClaim.
func createClaim(client *kubernetes.Clientset, params types.BuildParams, path string) error {
	claim, err := m8sclaim.Generate(m8sclaim.GenerateParams{
		Namespace:    params.Config.Namespace,
		Name:         utils.Machine(path),
		Annotations:  params.Annotations,
		StorageClass: params.Config.Cache.Type,
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate Ingress")
	}

	return skprclaim.Deploy(client, claim)
}

// Helper to create a Service.
func createService(client *kubernetes.Clientset, params types.BuildParams) error {
	svc, err := m8sservice.Generate(m8sservice.GenerateParams{
		Namespace:   params.Config.Namespace,
		Name:        params.Name,
		Port:        params.Config.Port,
		Annotations: params.Annotations,
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate Server")
	}

	return skprservice.Deploy(client, svc)
}

// Helper to create an Ingress.
func createIngress(client *kubernetes.Clientset, params types.BuildParams) error {
	ing, err := m8singress.Generate(m8singress.GenerateParams{
		Namespace:   params.Config.Namespace,
		Name:        params.Name,
		Domain:      params.Domain,
		Port:        params.Config.Port,
		Annotations: params.Annotations,
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate Ingress")
	}

	return skpringress.Deploy(client, ing)
}

// Helper to create a Pod.
func createPod(client *kubernetes.Clientset, params types.BuildParams) error {
	pod, err := m8spod.Generate(m8spod.GenerateParams{
		Namespace:       params.Config.Namespace,
		Name:            params.Name,
		Annotations:     params.Annotations,
		Repository:      params.Repository,
		Revision:        params.Revision,
		Services:        params.DockerCompose.Services,
		Caches:          params.Config.Cache.Paths,
		SecretDockerCfg: params.Config.Secrets.DockerCfg,
		SecretSSH:       params.Config.Secrets.SSH,
		Sidecar: sidecar.GenerateParams{
			User: params.Config.Auth.User,
			Pass: params.Config.Auth.Pass,
			Port: params.Config.Port,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate Pod")
	}

	return skprpod.Deploy(client, pod)
}
