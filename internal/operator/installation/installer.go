package installation

import (
	"context"

	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type installer struct {
	ClientSet    *kubernetes.Clientset
	Installation *model.Installation
	Logger       logr.Logger
}

func (i *installer) Run(ctx context.Context) (reconcile.Result, error) {
	if err := i.resolveVersion(ctx); err != nil {
		return reconcile.Result{}, err
	}

	if i.Installation.Status.DesiredVersion == nil {
		return reconcile.Result{RequeueAfter: defaultRetryInterval}, nil
	}

	r := &renderer{
		ClientSet:    i.ClientSet,
		Installation: i.Installation,
	}

	_, err := r.Render(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: defaultReconcileInterval}, nil
}

// func (i *installer) resolveInputs(ctx context.Context) error {
// 	configmaps := map[string]*corev1.ConfigMap{}
// 	secrets := map[string]*corev1.Secret{}

// 	for _, input := range i.Installation.Spec.Inputs {
// 		if input.ValueFrom != nil {
// 			if ref := input.ValueFrom.ConfigMapKeyRef; ref != nil {
// 				configs[ref.Name] = nil
// 			}
// 			if ref := input.ValueFrom.SecretKeyRef; ref != nil {
// 				secrets[ref.Name] = nil
// 			}
// 		}
// 	}

// 	for _, inputs := range i.Installation.Spec.InputsFrom {
// 		if ref := inputs.ConfigMapRef; ref != nil {
// 			configs[ref.Name] = nil
// 		}
// 		if ref := inputs.SecretRef; ref != nil {
// 			secrets[ref.Name] = nil
// 		}
// 	}

// 	for name := range configs {
// 		cm, err := i.
// 			ClientSet.
// 			CoreV1().
// 			ConfigMaps(i.Installation.Namespace).
// 			Watch()
// 			Get(ctx, name, metav1.GetOptions{})
// 		if err != nil && !apierrors.IsNotFound(err) {
// 			return err
// 		}

// 		configs[name] = cm
// 	}

// 	for name := range secrets {
// 		s, err := i.
// 			ClientSet.
// 			CoreV1().
// 			Secrets(i.Installation.Namespace).
// 			Get(ctx, name, metav1.GetOptions{})
// 	}

// }

// // func (r *Reconciler) reconcileJob(
// // 	ctx context.Context,
// // 	inst *model.Installation,
// // ) error {
// // 	k, err := kubernetes.NewForConfig(r.Manager.GetConfig())
// // 	if err != nil {
// // 		return err
// // 	}

// // 	name := fmt.Sprintf(
// // 		"cuerator-installation-%s",
// // 		inst.GetName(),
// // 	)

// // 	jobs := k.
// // 		BatchV1().
// // 		Jobs(inst.GetNamespace())

// // 	job, err := jobs.Get(ctx, name, metav1.GetOptions{})
// // 	if kerrors.IsNotFound(err) {
// // 		job, err = jobs.Create(
// // 			ctx,
// // 			buildJobSpec(name, inst),
// // 			metav1.CreateOptions{
// // 				FieldManager: "cuerator.io",
// // 			},
// // 		)
// // 	}
// // 	if err != nil {
// // 		return err
// // 	}

// // 	dapper.Print(job)

// // 	return nil
// // }

// // // func buildJobSpec(
// // // 	name string,
// // // 	inst *Installation,
// // // ) *batchv1.Job {
// // // 	return &batchv1.Job{
// // // 		TypeMeta: metav1.TypeMeta{
// // // 			APIVersion: "batch/v1",
// // // 			Kind:       "Job",
// // // 		},
// // // 		ObjectMeta: metav1.ObjectMeta{
// // // 			Namespace: inst.GetNamespace(),
// // // 			Name:      name,
// // // 			OwnerReferences: []metav1.OwnerReference{
// // // 				{
// // // 					APIVersion: inst.APIVersion,
// // // 					Kind:       inst.Kind,
// // // 					Name:       inst.GetName(),
// // // 					UID:        inst.GetUID(),
// // // 				},
// // // 			},
// // // 		},
// // // 		Spec: batchv1.JobSpec{
// // // 			Completions:  ptr[int32](1),
// // // 			BackoffLimit: ptr[int32](3),
// // // 			Template: corev1.PodTemplateSpec{
// // // 				Spec: corev1.PodSpec{
// // // 					RestartPolicy: corev1.RestartPolicyOnFailure,
// // // 					InitContainers: []corev1.Container{
// // // 						{
// // // 							Name:            "cuerator",
// // // 							Image:           fmt.Sprintf("ghcr.io/cuerator-io/cuerator:dev"),
// // // 							ImagePullPolicy: corev1.PullNever, // TODO
// // // 							Command:         []string{"/bin/cuerator-renderer"},
// // // 							Args:            []string{"init", "/mnt/share"},
// // // 							VolumeMounts: []corev1.VolumeMount{
// // // 								{
// // // 									Name:      "share",
// // // 									MountPath: "/mnt/share",
// // // 								},
// // // 							},
// // // 						},
// // // 					},
// // // 					Containers: []corev1.Container{
// // // 						{
// // // 							Name: "installation-templates",
// // // 							Image: fmt.Sprintf(
// // // 								"%s@%s",
// // // 								inst.Spec.Image,
// // // 								inst.Status.Tag.Digest,
// // // 							),
// // // 							Command: []string{"/.cuerator/bin/cuerator-renderer"},
// // // 							Args:    []string{"render"},
// // // 							VolumeMounts: []corev1.VolumeMount{
// // // 								{
// // // 									Name:      "share",
// // // 									MountPath: "/.cuerator",
// // // 								},
// // // 							},
// // // 						},
// // // 					},
// // // 					Volumes: []corev1.Volume{
// // // 						{
// // // 							Name: "share",
// // // 							VolumeSource: corev1.VolumeSource{
// // // 								EmptyDir: &corev1.EmptyDirVolumeSource{
// // // 									Medium:    corev1.StorageMediumMemory,
// // // 									SizeLimit: resource.NewScaledQuantity(32, resource.Mega),
// // // 								},
// // // 							},
// // // 						},
// // // 					},
// // // 				},
// // // 			},
// // // 		},
// // // 	}
// // // }

// // // func ptr[T any](v T) *T {
// // // 	return &v
// // // }

// import (
// 	"context"
// 	"fmt"

// 	batchv1 "k8s.io/api/batch/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	kerrors "k8s.io/apimachinery/pkg/api/errors"
// 	"k8s.io/apimachinery/pkg/api/resource"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/watch"
// )

// type Renderer struct {
// 	Client *kubernetes.Clientset
// }

// func (r *Renderer) Render(
// 	ctx context.Context,
// 	inst *Installation,
// ) ([]byte, error) {
// 	jobName := fmt.Sprintf(
// 		"cuerator-renderer-%s",
// 		inst.GetName(),
// 	)

// 	jobs := r.Client.
// 		BatchV1().
// 		Jobs(inst.GetNamespace())

// 	job, err := jobs.Get(ctx, jobName, metav1.GetOptions{})

// 	if err != nil {
// 		if !kerrors.IsNotFound(err) {
// 			return nil, fmt.Errorf("unable to lookup renderer job: %w", err)
// 		}

// 		job, err = r.createJob(ctx, inst)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	watch.Until(
// 		ctx,
// 	)

// 	return nil, nil
// }

// func (r *Renderer) createJob(
// 	ctx context.Context,
// 	inst *Installation,
// )

// func buildJobSpec(
// 	name string,
// 	inst *Installation,
// ) *batchv1.Job {
// 	return &batchv1.Job{
// 		TypeMeta: metav1.TypeMeta{
// 			APIVersion: "batch/v1",
// 			Kind:       "Job",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Namespace: inst.GetNamespace(),
// 			Name:      name,
// 			OwnerReferences: []metav1.OwnerReference{
// 				{
// 					APIVersion: inst.APIVersion,
// 					Kind:       inst.Kind,
// 					Name:       inst.GetName(),
// 					UID:        inst.GetUID(),
// 				},
// 			},
// 		},
// 		Spec: batchv1.JobSpec{
// 			Completions:  ptr[int32](1),
// 			BackoffLimit: ptr[int32](3),
// 			Template: corev1.PodTemplateSpec{
// 				Spec: corev1.PodSpec{
// 					RestartPolicy: corev1.RestartPolicyOnFailure,
// 					InitContainers: []corev1.Container{
// 						{
// 							Name:            "cuerator",
// 							Image:           fmt.Sprintf("ghcr.io/cuerator-io/cuerator:dev"),
// 							ImagePullPolicy: corev1.PullNever, // TODO
// 							Command:         []string{"/bin/cuerator-renderer"},
// 							Args:            []string{"init", "/mnt/share"},
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      "share",
// 									MountPath: "/mnt/share",
// 								},
// 							},
// 						},
// 					},
// 					Containers: []corev1.Container{
// 						{
// 							Name: "installation-templates",
// 							Image: fmt.Sprintf(
// 								"%s@%s",
// 								inst.Spec.Image,
// 								inst.Status.Tag.Digest,
// 							),
// 							Command: []string{"/.cuerator/bin/cuerator-renderer"},
// 							Args:    []string{"render"},
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      "share",
// 									MountPath: "/.cuerator",
// 								},
// 							},
// 						},
// 					},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: "share",
// 							VolumeSource: corev1.VolumeSource{
// 								EmptyDir: &corev1.EmptyDirVolumeSource{
// 									Medium:    corev1.StorageMediumMemory,
// 									SizeLimit: resource.NewScaledQuantity(32, resource.Mega),
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func ptr[T any](v T) *T {
// 	return &v
// }
