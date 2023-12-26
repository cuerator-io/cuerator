package di

import (
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/imbue"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var debug = ferrite.
	Bool("DEBUG", "enable verbose/debug logging").
	WithDefault(false).
	Required()

func init() {
	imbue.With0(
		Container,
		func(
			ctx imbue.Context,
		) (logr.Logger, error) {
			return zap.New(
				zap.UseDevMode(
					debug.Value(),
				),
			), nil
		},
	)
}
