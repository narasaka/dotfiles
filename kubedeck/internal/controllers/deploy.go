package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/kubedeck/kubedeck/internal/k8s"
	"github.com/kubedeck/kubedeck/internal/models"
)

type DeployController struct {
	deployments *models.DeploymentStore
	apps        *models.AppStore
	k8sClient   *k8s.Client
}

func NewDeployController(
	deployments *models.DeploymentStore,
	apps *models.AppStore,
	k8sClient *k8s.Client,
) *DeployController {
	return &DeployController{
		deployments: deployments,
		apps:        apps,
		k8sClient:   k8sClient,
	}
}

func (c *DeployController) Deploy(app *models.App, build *models.Build) (*models.Deployment, error) {
	ctx := context.Background()

	c.apps.UpdateStatus(app.ID, "deploying")

	deployName := fmt.Sprintf("kubedeck-%s", app.Name)

	// Create deployment record
	dep := &models.Deployment{
		AppID:             app.ID,
		BuildID:           build.ID,
		K8sDeploymentName: deployName,
		ReplicasDesired:   app.Replicas,
		ReplicasReady:     0,
		Status:            "pending",
	}

	created, err := c.deployments.Create(dep)
	if err != nil {
		return nil, fmt.Errorf("create deployment record: %w", err)
	}

	if c.k8sClient == nil {
		// Dev mode — simulate deployment
		c.deployments.UpdateStatus(created.ID, "running")
		c.deployments.UpdateReplicas(created.ID, app.Replicas)
		c.apps.UpdateStatus(app.ID, "running")
		return created, nil
	}

	envVars := k8s.ParseEnvVars(app.EnvVars)

	// Create/update K8s Deployment
	if err := c.k8sClient.CreateOrUpdateDeployment(ctx, k8s.DeployOpts{
		Name:      app.Name,
		Namespace: app.Namespace,
		AppID:     app.ID,
		Image:     build.ImageTag,
		Replicas:  int32(app.Replicas),
		Port:      int32(app.Port),
		EnvVars:   envVars,
	}); err != nil {
		c.deployments.UpdateStatus(created.ID, "failed")
		c.apps.UpdateStatus(app.ID, "failed")
		return nil, fmt.Errorf("create k8s deployment: %w", err)
	}

	// Create/update Service
	if err := c.k8sClient.CreateOrUpdateService(ctx, k8s.ServiceOpts{
		Name:      app.Name,
		Namespace: app.Namespace,
		AppID:     app.ID,
		Port:      int32(app.Port),
	}); err != nil {
		log.Printf("Warning: failed to create service for %s: %v", app.Name, err)
	}

	// Create/update Ingress if host is set
	if app.IngressHost != "" {
		if err := c.k8sClient.CreateOrUpdateIngress(ctx, k8s.IngressOpts{
			Name:      app.Name,
			Namespace: app.Namespace,
			AppID:     app.ID,
			Host:      app.IngressHost,
			Port:      int32(app.Port),
			TLS:       app.IngressTLS,
		}); err != nil {
			log.Printf("Warning: failed to create ingress for %s: %v", app.Name, err)
		}
	}

	c.deployments.UpdateStatus(created.ID, "rolling_out")

	// Watch rollout in background
	go c.watchRollout(ctx, app, created)

	return created, nil
}

func (c *DeployController) watchRollout(ctx context.Context, app *models.App, dep *models.Deployment) {
	deployName := fmt.Sprintf("kubedeck-%s", app.Name)

	for i := 0; i < 120; i++ { // 10 minutes max
		_, ready, err := c.k8sClient.GetDeploymentStatus(ctx, app.Namespace, deployName)
		if err != nil {
			log.Printf("Error checking deployment status: %v", err)
			continue
		}

		c.deployments.UpdateReplicas(dep.ID, int(ready))

		if int(ready) >= app.Replicas {
			c.deployments.UpdateStatus(dep.ID, "running")
			c.apps.UpdateStatus(app.ID, "running")
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	c.deployments.UpdateStatus(dep.ID, "failed")
	c.apps.UpdateStatus(app.ID, "failed")
}
