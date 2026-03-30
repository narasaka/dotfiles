package k8s

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KanikoBuildOpts struct {
	BuildID        string
	AppID          string
	GitURL         string
	GitBranch      string
	DockerfilePath string
	RegistryImage  string
	CommitSHA      string
	KanikoImage    string
}

func (c *Client) CreateKanikoJob(ctx context.Context, opts KanikoBuildOpts) (string, error) {
	shortID := opts.BuildID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	jobName := fmt.Sprintf("kubedeck-build-%s", shortID)
	ttl := int32(3600)
	backoffLimit := int32(0)

	destination := fmt.Sprintf("%s:%s", opts.RegistryImage, opts.CommitSHA)
	if len(opts.CommitSHA) > 7 {
		destination = fmt.Sprintf("%s:%s", opts.RegistryImage, opts.CommitSHA[:7])
	}

	gitContext := fmt.Sprintf("git://%s#refs/heads/%s", opts.GitURL, opts.GitBranch)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "kubedeck",
				"kubedeck/app-id":              opts.AppID,
				"kubedeck/build-id":            opts.BuildID,
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttl,
			BackoffLimit:            &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"kubedeck/build-id": opts.BuildID,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kaniko",
							Image: opts.KanikoImage,
							Args: []string{
								fmt.Sprintf("--dockerfile=%s", opts.DockerfilePath),
								fmt.Sprintf("--context=%s", gitContext),
								fmt.Sprintf("--destination=%s", destination),
								"--cache=true",
								"--snapshot-mode=redo",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "docker-config",
									MountPath: "/kaniko/.docker",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "docker-config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "kubedeck-registry-creds",
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
		return "", fmt.Errorf("create kaniko job: %w", err)
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
