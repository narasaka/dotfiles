package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceOpts struct {
	Name      string
	Namespace string
	AppID     string
	Port      int32
}

func (c *Client) CreateOrUpdateService(ctx context.Context, opts ServiceOpts) error {
	labels := map[string]string{
		"kubeploy/app-id":   opts.AppID,
		"kubeploy/app-name": opts.Name,
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("kubeploy-%s", opts.Name),
			Namespace: opts.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "kubeploy",
				"kubeploy/app-id":              opts.AppID,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       opts.Port,
					TargetPort: intstr.FromInt32(opts.Port),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	existing, err := c.Clientset.CoreV1().Services(opts.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = c.Clientset.CoreV1().Services(opts.Namespace).Create(ctx, svc, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	existing.Spec.Ports = svc.Spec.Ports
	existing.Spec.Selector = svc.Spec.Selector
	_, err = c.Clientset.CoreV1().Services(opts.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func (c *Client) DeleteService(ctx context.Context, namespace, name string) error {
	return c.Clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
