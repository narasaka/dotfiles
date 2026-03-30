package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubedeck/kubedeck/internal/k8s"
	"github.com/kubedeck/kubedeck/internal/models"
)

type BuildController struct {
	builds    *models.BuildStore
	apps      *models.AppStore
	settings  *models.SettingsStore
	k8sClient *k8s.Client
	deployer  *DeployController
}

func NewBuildController(
	builds *models.BuildStore,
	apps *models.AppStore,
	settings *models.SettingsStore,
	k8sClient *k8s.Client,
	deployer *DeployController,
) *BuildController {
	return &BuildController{
		builds:    builds,
		apps:      apps,
		settings:  settings,
		k8sClient: k8sClient,
		deployer:  deployer,
	}
}

func (c *BuildController) TriggerBuild(app *models.App, commitSHA, commitMessage, commitAuthor string) (*models.Build, error) {
	shortSHA := commitSHA
	if len(shortSHA) > 7 {
		shortSHA = shortSHA[:7]
	}

	imageTag := fmt.Sprintf("%s:%s", app.RegistryImage, shortSHA)

	build := &models.Build{
		AppID:         app.ID,
		CommitSHA:     commitSHA,
		CommitMessage: commitMessage,
		CommitAuthor:  commitAuthor,
		ImageTag:      imageTag,
		Status:        "pending",
	}

	created, err := c.builds.Create(build)
	if err != nil {
		return nil, fmt.Errorf("create build record: %w", err)
	}

	c.apps.UpdateStatus(app.ID, "building")
	c.apps.UpdateCurrentBuild(app.ID, created.ID)

	// Run build async
	go c.runBuild(app, created)

	return created, nil
}

func (c *BuildController) runBuild(app *models.App, build *models.Build) {
	ctx := context.Background()

	c.builds.SetStarted(build.ID)
	c.builds.AppendLogs(build.ID, fmt.Sprintf("=== Build started at %s ===\n", time.Now().Format(time.RFC3339)))

	if c.k8sClient == nil {
		c.builds.AppendLogs(build.ID, "K8s client not available (dev mode)\n")
		c.builds.AppendLogs(build.ID, "Simulating build...\n")
		time.Sleep(2 * time.Second)
		c.builds.AppendLogs(build.ID, "Build completed (simulated)\n")
		c.builds.SetFinished(build.ID, "success")
		c.apps.UpdateStatus(app.ID, "running")
		return
	}

	kanikoImage := "gcr.io/kaniko-project/executor:latest"
	if img, err := c.settings.Get("kaniko_image"); err == nil && img != "" {
		kanikoImage = img
	}

	jobName, err := c.k8sClient.CreateKanikoJob(ctx, k8s.KanikoBuildOpts{
		BuildID:        build.ID,
		AppID:          app.ID,
		GitURL:         app.GitURL,
		GitBranch:      app.GitBranch,
		DockerfilePath: app.DockerfilePath,
		RegistryImage:  app.RegistryImage,
		CommitSHA:      build.CommitSHA,
		KanikoImage:    kanikoImage,
	})
	if err != nil {
		c.builds.AppendLogs(build.ID, fmt.Sprintf("Failed to create Kaniko job: %v\n", err))
		c.builds.SetFinished(build.ID, "failed")
		c.apps.UpdateStatus(app.ID, "failed")
		return
	}

	c.builds.SetKanikoJobName(build.ID, jobName)
	c.builds.AppendLogs(build.ID, fmt.Sprintf("Created Kaniko job: %s\n", jobName))

	// Watch job status
	c.watchBuild(ctx, app, build, jobName)
}

func (c *BuildController) watchBuild(ctx context.Context, app *models.App, build *models.Build, jobName string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(30 * time.Minute)

	for {
		select {
		case <-timeout:
			c.builds.AppendLogs(build.ID, "Build timed out\n")
			c.builds.SetFinished(build.ID, "failed")
			c.apps.UpdateStatus(app.ID, "failed")
			return

		case <-ticker.C:
			status, err := c.k8sClient.GetJobStatus(ctx, c.k8sClient.Namespace, jobName)
			if err != nil {
				log.Printf("Error checking job status: %v", err)
				continue
			}

			// Try to get logs
			logs, err := c.k8sClient.GetJobLogs(ctx, c.k8sClient.Namespace, jobName)
			if err == nil && logs != "" {
				c.builds.AppendLogs(build.ID, logs)
			}

			switch status {
			case "success":
				c.builds.AppendLogs(build.ID, fmt.Sprintf("\n=== Build succeeded at %s ===\n", time.Now().Format(time.RFC3339)))
				c.builds.SetFinished(build.ID, "success")

				// Trigger deploy
				if c.deployer != nil {
					if _, err := c.deployer.Deploy(app, build); err != nil {
						log.Printf("Deploy failed for app %s: %v", app.ID, err)
						c.apps.UpdateStatus(app.ID, "failed")
					}
				}
				return

			case "failed":
				c.builds.AppendLogs(build.ID, fmt.Sprintf("\n=== Build failed at %s ===\n", time.Now().Format(time.RFC3339)))
				c.builds.SetFinished(build.ID, "failed")
				c.apps.UpdateStatus(app.ID, "failed")
				return
			}
		}
	}
}

func (c *BuildController) CancelBuild(build *models.Build) error {
	c.builds.SetFinished(build.ID, "cancelled")
	c.builds.AppendLogs(build.ID, "Build cancelled by user\n")

	if c.k8sClient != nil && build.KanikoJobName != "" {
		ctx := context.Background()
		c.k8sClient.DeleteJob(ctx, c.k8sClient.Namespace, build.KanikoJobName)
	}

	return nil
}
