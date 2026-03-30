package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func ListOptionsForManagedBy() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=kubedeck",
	}
}
