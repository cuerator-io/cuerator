package installation

import (
	"context"
	"time"

	"github.com/cuerator-io/cuerator/internal/operator"
	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	"github.com/dogmatiq/dyad"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultReconcileInterval = 5 * time.Minute
	defaultRetryInterval     = 15 * time.Second
	finalizerName            = operator.GroupName + "/uninstall"
)

// reconciler is an implementation of [reconcile.Reconciler] that manages
// [model.Installation] resources.
type reconciler struct {
	Manager manager.Manager
	Logger  logr.Logger
}

func (r *reconciler) Reconcile(
	ctx context.Context,
	req reconcile.Request,
) (reconcile.Result, error) {
	logger := r.Logger.WithName(req.Namespace).WithName(req.Name)

	inst := &model.Installation{}
	if err := r.Manager.GetClient().Get(ctx, req.NamespacedName, inst); err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return reconcile.Result{}, err
		}

		logger.Info("resource deleted")
		return reconcile.Result{}, nil
	}

	logger.Info(
		"loaded resource",
		"generation", inst.Generation,
		"version", inst.ResourceVersion,
	)

	if logger.Enabled() {
		orig := dyad.Clone(inst)

		defer func() {
			if diff := cmp.Diff(orig.Status, inst.Status); diff != "" {
				logger.Info("status updated\n" + diff)
			}
		}()
	}

	if !inst.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, r.uninstall(ctx, inst, logger)
	}

	res, err := r.install(ctx, inst, logger)
	if err != nil {
		return reconcile.Result{}, err
	}

	if res.Requeue || res.RequeueAfter > 0 {
		logger.Info(
			"requeueing reconciliation",
			"delay", res.RequeueAfter,
		)
	}

	return res, nil
}

func (r *reconciler) install(
	ctx context.Context,
	inst *model.Installation,
	logger logr.Logger,
) (reconcile.Result, error) {
	c := r.Manager.GetClient()

	if controllerutil.AddFinalizer(inst, finalizerName) {
		if err := c.Update(ctx, inst); err != nil {
			return reconcile.Result{}, err
		}
		logger.Info(
			"added finalizer",
			"finalizer", finalizerName,
		)
	}

	clientSet, err := kubernetes.NewForConfig(r.Manager.GetConfig())
	if err != nil {
		return reconcile.Result{}, err
	}

	op := &installer{
		Installation: inst,
		ClientSet:    clientSet,
		Logger:       logger,
	}

	res, err := op.Run(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	return res, c.Status().Update(ctx, inst)
}

func (r *reconciler) uninstall(
	ctx context.Context,
	inst *model.Installation,
	logger logr.Logger,
) error {
	op := &uninstaller{
		Installation: inst,
		Logger:       logger,
	}

	if err := op.Run(ctx); err != nil {
		return err
	}

	c := r.Manager.GetClient()
	if err := c.Status().Update(ctx, inst); err != nil {
		return err
	}

	if controllerutil.RemoveFinalizer(inst, finalizerName) {
		if err := c.Update(ctx, inst); err != nil {
			return err
		}
		logger.Info(
			"removed finalizer",
			"finalizer", finalizerName,
		)
	}

	return nil
}
