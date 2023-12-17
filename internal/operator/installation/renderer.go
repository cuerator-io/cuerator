package installation

import (
	"context"
	"fmt"
	"os"

	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	"github.com/dogmatiq/dapper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type renderer struct {
	ClientSet    *kubernetes.Clientset
	Installation *model.Installation
}

func (r *renderer) Render(ctx context.Context) ([]byte, error) {
	name := fmt.Sprintf(
		"cuerator-renderer-%s",
		r.Installation.GetName(),
	)

	jobs := r.ClientSet.
		BatchV1().
		Jobs(r.Installation.GetNamespace())

	job, err := jobs.Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("unable to find renderer job: %w", err)
		}

		job, err = jobs.Create(
			ctx,
			r.buildJob(name),
			metav1.CreateOptions{
				FieldManager: "cuerator.io",
			},
		)
		if err != nil {
			return nil, err
		}
	}

	dapper.Print(job)

	return nil, nil
}

func (r *renderer) buildJob(name string) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.Installation.Namespace,
			Name:      name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: r.Installation.APIVersion,
					Kind:       r.Installation.Kind,
					Name:       r.Installation.Name,
					UID:        r.Installation.UID,
				},
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					InitContainers: []corev1.Container{
						{
							Name:            "copy-self",
							Image:           fmt.Sprintf("ghcr.io/cuerator-io/cuerator:dev"),
							ImagePullPolicy: corev1.PullNever, // TODO
							Command:         []string{"/bin/cuerator"},
							Args: []string{
								"copy-self",
								"/.cuerator/bin/cuerator",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "share",
									MountPath: "/.cuerator/bin",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "render",
							Image: fmt.Sprintf(
								"%s@%s",
								r.Installation.Status.DesiredVersion.Image,
								r.Installation.Status.DesiredVersion.Digest,
							),
							Command: []string{"/.cuerator/bin/cuerator"},
							Args: []string{
								"render",
								"/.cuerator",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "share",
									MountPath: "/.cuerator/bin",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "share",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium:    corev1.StorageMediumMemory,
									SizeLimit: resource.NewQuantity(volumeSize, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}

// volumeSize is the size of the volume used to share data between the Cuerator
// process and the container that houses the Cuerator collection.
var volumeSize int64 = 10 * 1024 * 1024 // 10 MiB

// Add the size of the Cuerator executable to the volume size.
func init() {
	filename, err := os.Executable()
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}

	volumeSize += info.Size()
}
