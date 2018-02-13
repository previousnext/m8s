package cronjob

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"github.com/pkg/errors"
)

// Deploy will create the CronJob and fallback to updating if it exists.
func Deploy(client *kubernetes.Clientset, cronjob *batchv1beta1.CronJob) error {
	_, err := client.BatchV1beta1().CronJobs(cronjob.ObjectMeta.Namespace).Create(cronjob)
	if kerrors.IsAlreadyExists(err) {
		_, err := client.BatchV1beta1().CronJobs(cronjob.ObjectMeta.Namespace).Update(cronjob)
		if err != nil {
			return errors.Wrap(err, "failed to update")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
