package k8s

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildOpts struct {
	BuildID        string
	AppID          string
	GitURL         string
	GitBranch      string
	DockerfilePath string
	RegistryImage  string
	CommitSHA      string
	BuildKitAddr   string
}

func (c *Client) CreateBuildJob(ctx context.Context, opts BuildOpts) (string, error) {
	shortID := opts.BuildID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	jobName := fmt.Sprintf("kubeploy-build-%s", shortID)
	ttl := int32(3600)
	backoffLimit := int32(0)

	destination := fmt.Sprintf("%s:%s", opts.RegistryImage, opts.CommitSHA)
	if len(opts.CommitSHA) > 7 {
		destination = fmt.Sprintf("%s:%s", opts.RegistryImage, opts.CommitSHA[:7])
	}

	buildkitAddr := opts.BuildKitAddr
	if buildkitAddr == "" {
		buildkitAddr = "tcp://kubeploy-buildkitd:1234"
	}

	gitRef := fmt.Sprintf("refs/heads/%s", opts.GitBranch)
	dockerfile := opts.DockerfilePath
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "kubeploy",
				"kubeploy/app-id":              opts.AppID,
				"kubeploy/build-id":            opts.BuildID,
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttl,
			BackoffLimit:            &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"kubeploy/build-id": opts.BuildID,
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "git-clone",
							Image: "alpine/git:latest",
							Command: []string{"sh", "-c"},
							Args: []string{
								fmt.Sprintf("git clone --depth 1 --branch %s %s /workspace", opts.GitBranch, opts.GitURL),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "workspace",
									MountPath: "/workspace",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "buildkit",
							Image: "moby/buildkit:v0.21.1",
							Command: []string{"buildctl"},
							Args: []string{
								"--addr", buildkitAddr,
								"build",
								"--frontend", "dockerfile.v0",
								"--local", "context=/workspace",
								"--local", "dockerfile=/workspace",
								"--opt", fmt.Sprintf("filename=%s", dockerfile),
								"--output", fmt.Sprintf("type=image,name=%s,push=true", destination),
							},
							Env: []corev1.EnvVar{
								{
									Name:  "DOCKER_CONFIG",
									Value: "/docker-config",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "workspace",
									MountPath: "/workspace",
									ReadOnly:  true,
								},
								{
									Name:      "docker-config",
									MountPath: "/docker-config",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "workspace",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "docker-config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "kubeploy-registry-creds",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	created, err := c.Clientset.BatchV1().Jobs(c.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("create build job: %w", err)
	}

	return created.Name, nil
}

func (c *Client) DeleteJob(ctx context.Context, namespace, name string) error {
	propagation := metav1.DeletePropagationBackground
	return c.Clientset.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
}

func (c *Client) GetJobStatus(ctx context.Context, namespace, name string) (string, error) {
	job, err := c.Clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if job.Status.Succeeded > 0 {
		return "success", nil
	}
	if job.Status.Failed > 0 {
		return "failed", nil
	}
	if job.Status.Active > 0 {
		return "building", nil
	}
	return "pending", nil
}
