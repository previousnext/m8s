package podutils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	pb "github.com/previousnext/m8s/pb"
)

const buildExecuteFailed = "BuildExecuteFailed"

// Tail tails the logs for a build.
// @todo, Make the stream more generic.
func Tail(ctx context.Context, stream pb.M8S_CreateServer, client *kubernetes.Clientset, namespace, name string) error {
	pods := client.CoreV1().Pods(namespace)

	watcher := podWatcher{
		pods: pods,
		name: name,
	}
	if err := watcher.start(ctx); err != nil {
		return fmt.Errorf("watching pod: %v", err)
	}

	pod, err := watcher.waitForPod(ctx, func(p *v1.Pod) bool {
		return len(p.Status.InitContainerStatuses) > 0
	})
	if err != nil {
		return err
	}

	for i, container := range pod.Status.InitContainerStatuses {
		pod, err := watcher.waitForPod(ctx, func(p *v1.Pod) bool {
			waiting := p.Status.InitContainerStatuses[i].State.Waiting
			if waiting == nil {
				return true
			}

			if waiting.Message != "" {
				stream.Send(&pb.CreateResponse{
					Message: fmt.Sprintf("[%s] %s", container.Name, waiting.Message),
				})
			}

			return false
		})
		if err != nil {
			return fmt.Errorf("waiting for container: %v", err)
		}

		container := pod.Status.InitContainerStatuses[i]
		followContainer := container.State.Terminated == nil
		if err := printContainerLogs(ctx, stream, pods, name, container.Name, followContainer); err != nil {
			return fmt.Errorf("printing logs: %v", err)
		}

		pod, err = watcher.waitForPod(ctx, func(p *v1.Pod) bool {
			// Ensure that our container status exists.
			if len(p.Status.InitContainerStatuses) < i {
				return false
			}

			return p.Status.InitContainerStatuses[i].State.Terminated != nil
		})
		if err != nil {
			return fmt.Errorf("waiting for container termination: %v", err)
		}

		container = pod.Status.InitContainerStatuses[i]
		terminated := container.State.Terminated
		if terminated.ExitCode != 0 {
			message := "Build Failed"
			if terminated.Message != "" {
				message += ": " + terminated.Message
			}

			return stream.Send(&pb.CreateResponse{
				Message: fmt.Sprintf("[%s] %s", container.Name, message),
			})
		}
	}

	return nil
}

type podWatcher struct {
	pods corev1.PodInterface
	name string

	versions chan *v1.Pod
	last     *v1.Pod
}

func (w *podWatcher) start(ctx context.Context) error {
	w.versions = make(chan *v1.Pod, 100)

	watcher, err := w.pods.Watch(metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", w.name).String(),
	})
	if err != nil {
		return fmt.Errorf("watching pod: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				watcher.Stop()
				return
			case evt, ok := <-watcher.ResultChan():
				if !ok {
					continue
				}

				if evt.Object == nil {
					continue
				}

				w.versions <- evt.Object.(*v1.Pod)
			}
		}
	}()

	return nil
}

func (w *podWatcher) waitForPod(ctx context.Context, predicate func(pod *v1.Pod) bool) (*v1.Pod, error) {
	if w.last != nil && predicate(w.last) {
		return w.last, nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("watch cancelled")
		case w.last = <-w.versions:
			if predicate(w.last) {
				return w.last, nil
			}
		}
	}
}

func printContainerLogs(ctx context.Context, stream pb.M8S_CreateServer, pods corev1.PodExpansion, podName, containerName string, follow bool) error {
	rc, err := pods.GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	}).Stream()
	if err != nil {
		return err
	}
	defer rc.Close()

	return streamLogs(ctx, stream, containerName, rc)
}

func streamLogs(ctx context.Context, stream pb.M8S_CreateServer, containerName string, rc io.Reader) error {
	prefix := fmt.Sprintf("[%s]", containerName)

	r := bufio.NewReader(rc)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			if len(line) > 0 {
				err := stream.Send(&pb.CreateResponse{
					Message: fmt.Sprintf("%s %s", prefix, line),
				})
				if err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			return err
		}

		err = stream.Send(&pb.CreateResponse{
			Message: fmt.Sprintf("%s %s", prefix, line),
		})
		if err != nil {
			return err
		}
	}
}
