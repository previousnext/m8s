package deployment

import (
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// TimedOutReason is for checking the status for the deploy
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/util/deployment_util.go#L93
	TimedOutReason = "ProgressDeadlineExceeded"
)

// Deploy will create/update the deployment and wait for it to finish rolling out.
func Deploy(client *kubernetes.Clientset, deployment *appsv1.Deployment) error {
	dply, err := client.AppsV1().Deployments(deployment.ObjectMeta.Namespace).Create(deployment)
	if kerrors.IsAlreadyExists(err) {
		dply, err = client.AppsV1().Deployments(deployment.ObjectMeta.Namespace).Update(deployment)
		if err != nil {
			return errors.Wrap(err, "failed to update Deployment")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to create Deployment")
	}

	// Wait for the rollout to complete.
	limiter := time.Tick(time.Second * 15)

	for {
		<-limiter

		// Query the deployment to get a status update.
		dply, err := client.AppsV1().Deployments(deployment.ObjectMeta.Namespace).Get(dply.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed while waiting to finish")
		}

		if dply.Generation <= dply.Status.ObservedGeneration {
			cond := condition(dply.Status, appsv1.DeploymentProgressing)

			if cond != nil && cond.Reason == TimedOutReason {
				return errors.New("failed because of exceeded deadline")
			}

			if dply.Status.UpdatedReplicas < *dply.Spec.Replicas {
				continue
			}

			if dply.Status.Replicas > dply.Status.UpdatedReplicas {
				continue
			}

			if dply.Status.AvailableReplicas < dply.Status.UpdatedReplicas {
				continue
			}

			return nil
		}
	}

	return errors.New("failed to with reason unknown")
}

// condition returns the condition with the provided type.
// Borrowed from: https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/util/deployment_util.go#L117
func condition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]

		if c.Type == condType {
			return &c
		}
	}
	return nil
}
