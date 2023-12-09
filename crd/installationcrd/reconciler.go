package installationcrd

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles the state of [Installation] resources.
type Reconciler struct {
	Manager manager.Manager
	Client  client.Client
	Logger  logr.Logger
}

// Reconcile performs a full reconciliation for the object referred to by the
// request.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	inst := &Installation{}
	if err := r.Client.Get(ctx, req.NamespacedName, inst); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// if err := r.resolveVersion(ctx, inst); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// if err := r.reconcileJob(ctx, inst); err != nil {
	// 	return reconcile.Result{}, err
	// }

	return reconcile.Result{
		RequeueAfter: 5 * time.Minute,
	}, nil
}

func (r *Reconciler) setConditionOnReturn(
	ctx context.Context,
	inst *Installation,
	cond *metav1.Condition,
	err *error,
) {
	if cond.Status != metav1.ConditionUnknown && cond.Reason == "" {
		panic("the condition must have a reason")
	}

	if *err == nil {
		cond.Status = metav1.ConditionTrue
	} else {
		cond.Message = (*err).Error()
		cond.Status = metav1.ConditionFalse
	}

	cond.ObservedGeneration = inst.Generation
	inst.Status.Conditions.Merge(*cond)

	*err = errors.Join(
		*err,
		r.Client.Status().Update(ctx, inst),
	)
}
