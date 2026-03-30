package k8s

import (
	"bufio"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) StreamPodLogs(ctx context.Context, namespace, labelSelector string) (<-chan string, error) {
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	ch := make(chan string, 100)

	for _, pod := range pods.Items {
		podName := pod.Name
		go func() {
			tailLines := int64(100)
			req := c.Clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
				Follow:    true,
				TailLines: &tailLines,
			})

			stream, err := req.Stream(ctx)
			if err != nil {
				ch <- fmt.Sprintf("[%s] error: %v\n", podName, err)
				return
			}
			defer stream.Close()

			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				case ch <- fmt.Sprintf("[%s] %s\n", podName, scanner.Text()):
				}
			}
		}()
	}

	return ch, nil
}

func (c *Client) GetJobLogs(ctx context.Context, namespace, jobName string) (string, error) {
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return "", fmt.Errorf("list job pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found for job %s", jobName)
	}

	req := c.Clientset.CoreV1().Pods(namespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{})
	result, err := req.Do(ctx).Raw()
	if err != nil {
		return "", fmt.Errorf("get logs: %w", err)
	}

	return string(result), nil
}
