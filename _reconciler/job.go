package reconciler

import (
	"context"
	"fmt"

	"github.com/cuerator-io/cuerator/crd/installationcrd"
	"github.com/dogmatiq/dapper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (r *Reconciler) reconcileJob(
	ctx context.Context,
	inst *installationcrd.Installation,
) error {
	k, err := kubernetes.NewForConfig(r.Manager.GetConfig())
	if err != nil {
		return err
	}

	name := fmt.Sprintf(
		"cuerator-installation-%s",
		inst.GetName(),
	)

	jobs := k.
		BatchV1().
		Jobs(inst.GetNamespace())

	job, err := jobs.Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		job, err = jobs.Create(
			ctx,
			buildJobSpec(name, inst),
			metav1.CreateOptions{
				FieldManager: "cuerator.io",
			},
		)
	}
	if err != nil {
		return err
	}

	dapper.Print(job)

	return nil
}

func buildJobSpec(
	name string,
	inst *installationcrd.Installation,
) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: inst.GetNamespace(),
			Name:      name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: inst.APIVersion,
					Kind:       inst.Kind,
					Name:       inst.GetName(),
					UID:        inst.GetUID(),
				},
			},
		},
		Spec: batchv1.JobSpec{
			Completions:  ptr[int32](1),
			BackoffLimit: ptr[int32](3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					InitContainers: []corev1.Container{
						{
							Name:            "cuerator",
							Image:           fmt.Sprintf("ghcr.io/cuerator-io/cuerator:dev"),
							ImagePullPolicy: corev1.PullNever, // TODO
							Command:         []string{"/bin/cuerator-job"},
							Args:            []string{"init", "/mnt/share"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "share",
									MountPath: "/mnt/share",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "installation-templates",
							Image: fmt.Sprintf(
								"%s@%s",
								inst.Spec.Image,
								inst.Status.Tag.Digest,
							),
							Command: []string{"/.cuerator/bin/cuerator-job"},
							Args:    []string{"render"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "share",
									MountPath: "/.cuerator",
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
									SizeLimit: resource.NewScaledQuantity(32, resource.Mega),
								},
							},
						},
					},
				},
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
