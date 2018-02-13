package provisioner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"fmt"
	skprtesting "github.com/previousnext/skpr/k8stest/generate"
	"github.com/previousnext/skpr/k8stest/hash"
)

func TestRun(t *testing.T) {
	client, err := kubernetes.NewForConfig(&rest.Config{
		// This is our test "mock" cluster.
		Host: "http://localhost:8080",
	})

	var (
		name = hash.Generate(10)

		// We will write all our output to this so we can inspect in our test suite.
		buffer bytes.Buffer

		namespace  = skprtesting.Namespace()
		configmap  = skprtesting.ConfigMap(namespace.ObjectMeta.Name, name)
		deployment = skprtesting.Deployment(namespace.ObjectMeta.Name, name)
		service    = skprtesting.Service(namespace.ObjectMeta.Name, name)
		ingress    = skprtesting.Ingress(namespace.ObjectMeta.Name, name)
		claims     = skprtesting.PVCs(namespace.ObjectMeta.Name, "public", "private")
		cronjob    = skprtesting.CronJob(namespace.ObjectMeta.Name, name)
	)

	params := RunParams{
		Namespaces:             []*corev1.Namespace{namespace},
		ConfigMaps:             []*corev1.ConfigMap{configmap},
		Deployments:            []*appsv1.Deployment{deployment},
		Services:               []*corev1.Service{service},
		Ingresses:              []*extensionsv1beta1.Ingress{ingress},
		PersistentVolumeClaims: claims,
		CronJobs:               []*batchv1beta1.CronJob{cronjob},
	}

	err = Run(client, params, &buffer)
	assert.Nil(t, err)

	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: Namespace: %s", namespace.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: ConfigMap: %s", configmap.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), "Provisioning: PersistentVolumeClaim: public"))
	assert.True(t, strings.Contains(buffer.String(), "Provisioning: PersistentVolumeClaim: private"))
	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: Deployment: %s", deployment.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: Service: %s", service.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: Ingress: %s", ingress.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), fmt.Sprintf("Provisioning: CronJob: %s", cronjob.ObjectMeta.Name)))
	assert.True(t, strings.Contains(buffer.String(), "Complete!"))
}
