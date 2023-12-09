package main

import (
	"context"
	"log"

	"github.com/cuerator-io/cuerator/crd/installationcrd"
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/imbue"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	controller "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var container = imbue.New()

func main() {
	ferrite.Init()

	ctx := controller.SetupSignalHandler()

	if err := imbue.Invoke2(
		ctx,
		container,
		func(
			ctx context.Context,
			m manager.Manager,
			r *installationcrd.Reconciler,
		) error {
			if err := builder.
				ControllerManagedBy(m).
				For(&installationcrd.Installation{}).
				WithEventFilter(predicate.GenerationChangedPredicate{}).
				Complete(r); err != nil {
				return err
			}

			return m.Start(ctx)
		},
	); err != nil {
		log.Fatal(err)
	}
}
