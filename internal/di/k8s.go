package di

import (
	"github.com/dogmatiq/imbue"
	"github.com/go-logr/logr"
	controller "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	imbue.With1(
		Container,
		func(
			ctx imbue.Context,
			logger logr.Logger,
		) (manager.Manager, error) {
			cfg, err := controller.GetConfig()
			if err != nil {
				return nil, err
			}

			return controller.NewManager(
				cfg,
				controller.Options{
					Logger: logger,
				},
			)
		},
	)
}
