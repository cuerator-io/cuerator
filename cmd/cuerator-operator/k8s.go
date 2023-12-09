package main

import (
	"github.com/cuerator-io/cuerator/crd"
	"github.com/cuerator-io/cuerator/crd/installationcrd"
	"github.com/dogmatiq/imbue"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	controller "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

func init() {
	imbue.With1(
		container,
		func(
			ctx imbue.Context,
			l imbue.ByName[systemLogger, logr.Logger],
		) (manager.Manager, error) {
			cfg, err := controller.GetConfig()
			if err != nil {
				return nil, err
			}

			return controller.NewManager(
				cfg,
				controller.Options{
					Logger: l.Value(),
				},
			)
		},
	)

	imbue.With2(
		container,
		func(
			ctx imbue.Context,
			m manager.Manager,
			l imbue.ByName[verboseLogger, logr.Logger],
		) (*installationcrd.Reconciler, error) {
			return &installationcrd.Reconciler{
				Manager: m,
				Client:  m.GetClient(),
				Logger:  l.Value(),
			}, nil
		},
	)

	imbue.Decorate0(
		container,
		func(
			ctx imbue.Context,
			m manager.Manager,
		) (manager.Manager, error) {
			b := &scheme.Builder{
				GroupVersion: schema.GroupVersion{
					Group:   crd.GroupName,
					Version: crd.Version,
				},
			}

			b.Register(
				&installationcrd.Installation{},
				&installationcrd.InstallationList{},
			)

			if err := b.AddToScheme(m.GetScheme()); err != nil {
				return nil, err
			}

			return m, nil
		},
	)
}
