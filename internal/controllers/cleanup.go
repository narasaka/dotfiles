package controllers

import (
	"context"
	"log"
	"time"

	"github.com/kubeploy/kubeploy/internal/k8s"
)

type CleanupController struct {
	k8sClient *k8s.Client
}

func NewCleanupController(k8sClient *k8s.Client) *CleanupController {
	return &CleanupController{k8sClient: k8sClient}
}

func (c *CleanupController) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *CleanupController) cleanup(ctx context.Context) {
	if c.k8sClient == nil {
		return
	}

	log.Println("Running cleanup of old build jobs...")

	jobs, err := c.k8sClient.Clientset.BatchV1().Jobs(c.k8sClient.Namespace).List(ctx, k8s.ListOptionsForManagedBy())
	if err != nil {
		log.Printf("Cleanup: failed to list jobs: %v", err)
		return
	}

	cutoff := time.Now().Add(-24 * time.Hour)
	for _, job := range jobs.Items {
		if job.Status.CompletionTime != nil && job.Status.CompletionTime.Time.Before(cutoff) {
			if err := c.k8sClient.DeleteJob(ctx, c.k8sClient.Namespace, job.Name); err != nil {
				log.Printf("Cleanup: failed to delete job %s: %v", job.Name, err)
			} else {
				log.Printf("Cleanup: deleted old job %s", job.Name)
			}
		}
	}
}
