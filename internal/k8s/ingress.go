package k8s

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressOpts struct {
	Name      string
	Namespace string
	AppID     string
	Host      string
	Port      int32
	TLS       bool
}

func (c *Client) CreateOrUpdateIngress(ctx context.Context, opts IngressOpts) error {
	if opts.Host == "" {
		return nil // No ingress needed
	}

	pathType := networkingv1.PathTypePrefix
	ingressName := fmt.Sprintf("kubeploy-%s", opts.Name)
	svcName := fmt.Sprintf("kubeploy-%s", opts.Name)

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: opts.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "kubeploy",
				"kubeploy/app-id":              opts.AppID,
			},
			Annotations: map[string]string{},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: opts.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: svcName,
											Port: networkingv1.ServiceBackendPort{
												Number: opts.Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if opts.TLS {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts:      []string{opts.Host},
				SecretName: fmt.Sprintf("kubeploy-%s-tls", opts.Name),
			},
		}
		ingress.Annotations["cert-manager.io/cluster-issuer"] = "letsencrypt-prod"
	}

	existing, err := c.Clientset.NetworkingV1().Ingresses(opts.Namespace).Get(ctx, ingressName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = c.Clientset.NetworkingV1().Ingresses(opts.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	existing.Spec = ingress.Spec
	existing.Annotations = ingress.Annotations
	_, err = c.Clientset.NetworkingV1().Ingresses(opts.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func (c *Client) DeleteIngress(ctx context.Context, namespace, name string) error {
	return c.Clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
