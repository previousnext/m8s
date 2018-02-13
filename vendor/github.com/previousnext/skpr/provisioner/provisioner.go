package provisioner

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/skpr/utils/k8s/configmap"
	"github.com/previousnext/skpr/utils/k8s/cronjob"
	"github.com/previousnext/skpr/utils/k8s/deployment"
	"github.com/previousnext/skpr/utils/k8s/ingress"
	"github.com/previousnext/skpr/utils/k8s/namespace"
	"github.com/previousnext/skpr/utils/k8s/pvc"
	"github.com/previousnext/skpr/utils/k8s/service"
)

// RunParams are passed into the Run function.
type RunParams struct {
	Namespaces             []*corev1.Namespace
	ConfigMaps             []*corev1.ConfigMap
	PersistentVolumeClaims []*corev1.PersistentVolumeClaim
	Deployments            []*appsv1.Deployment
	Services               []*corev1.Service
	Ingresses              []*extensionsv1beta1.Ingress
	CronJobs               []*batchv1beta1.CronJob
}

// Run will start the provisioning process.
func Run(client *kubernetes.Clientset, params RunParams, w io.Writer) error {
	for _, ns := range params.Namespaces {
		fmt.Fprintln(w, "Provisioning: Namespace:", ns.ObjectMeta.Name)

		err := namespace.Deploy(client, ns)
		if err != nil {
			return errors.Wrap(err, "failed to provision: Namespace")
		}
	}

	for _, cfg := range params.ConfigMaps {
		fmt.Fprintln(w, "Provisioning: ConfigMap:", cfg.ObjectMeta.Name)

		err := configmap.Deploy(client, cfg)
		if err != nil {
			return errors.Wrap(err, "failed to provision: ConfigMap")
		}
	}

	for _, claim := range params.PersistentVolumeClaims {
		fmt.Fprintln(w, "Provisioning: PersistentVolumeClaim:", claim.ObjectMeta.Name)

		err := pvc.Deploy(client, claim)
		if err != nil {
			return errors.Wrap(err, "failed to provision: PersistentVolumeClaim")
		}
	}

	for _, dply := range params.Deployments {
		fmt.Fprintln(w, "Provisioning: Deployment:", dply.ObjectMeta.Name)

		err := deployment.Deploy(client, dply)
		if err != nil {
			return errors.Wrap(err, "failed to provision: Deployment")
		}
	}

	for _, svc := range params.Services {
		fmt.Fprintln(w, "Provisioning: Service:", svc.ObjectMeta.Name)

		err := service.Deploy(client, svc)
		if err != nil {
			return errors.Wrap(err, "failed to provision: Service")
		}
	}

	for _, ing := range params.Ingresses {
		fmt.Fprintln(w, "Provisioning: Ingress:", ing.ObjectMeta.Name)

		err := ingress.Deploy(client, ing)
		if err != nil {
			return errors.Wrap(err, "failed to provision: Ingress")
		}
	}

	for _, cron := range params.CronJobs {
		fmt.Fprintln(w, "Provisioning: CronJob:", cron.ObjectMeta.Name)

		err := cronjob.Deploy(client, cron)
		if err != nil {
			return errors.Wrap(err, "failed to provision: CronJob")
		}
	}

	fmt.Fprintln(w, "Complete!")

	return nil
}
