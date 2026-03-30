package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset *kubernetes.Clientset
	Config    *rest.Config
	Namespace string
}

func NewClient(namespace string) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig for local dev
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		if envKC := os.Getenv("KUBECONFIG"); envKC != "" {
			kubeconfig = envKC
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create k8s config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %w", err)
	}

	return &Client{
		Clientset: clientset,
		Config:    config,
		Namespace: namespace,
	}, nil
}
