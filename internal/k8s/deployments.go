package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeployOpts struct {
	Name      string
	Namespace string
	AppID     string
	Image     string
	Replicas  int32
	Port      int32
	EnvVars   map[string]string
}

func (c *Client) CreateOrUpdateDeployment(ctx context.Context, opts DeployOpts) error {
	labels := map[string]string{
		"app.kubernetes.io/managed-by": "kubeploy",
		"kubeploy/app-id":              opts.AppID,
		"kubeploy/app-name":            opts.Name,
	}

	var envList []corev1.EnvVar
	for k, v := range opts.EnvVars {
		envList = append(envList, corev1.EnvVar{Name: k, Value: v})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("kubeploy-%s", opts.Name),
			Namespace: opts.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &opts.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  opts.Name,
							Image: opts.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: opts.Port},
							},
							Env: envList,
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt32(opts.Port),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}

	existing, err := c.Clientset.AppsV1().Deployments(opts.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = c.Clientset.AppsV1().Deployments(opts.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	existing.Spec = deployment.Spec
	_, err = c.Clientset.AppsV1().Deployments(opts.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func (c *Client) DeleteDeployment(ctx context.Context, namespace, name string) error {
	return c.Clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *Client) GetDeploymentStatus(ctx context.Context, namespace, name string) (int32, int32, error) {
	dep, err := c.Clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, 0, err
	}
	return dep.Status.Replicas, dep.Status.ReadyReplicas, nil
}

func ParseEnvVars(envJSON string) map[string]string {
	result := make(map[string]string)
	json.Unmarshal([]byte(envJSON), &result)
	return result
}
