package di

import (
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/imbue"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type (
	systemLogger  imbue.Name[logr.Logger]
	verboseLogger imbue.Name[logr.Logger]
)

var debug = ferrite.
	Bool("DEBUG", "Enable debug logging.").
	WithDefault(false).
	Required()

func init() {
	imbue.With0Named[systemLogger](
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

	imbue.With1Named[verboseLogger](
		Container,
		func(
			ctx imbue.Context,
			l imbue.ByName[systemLogger, logr.Logger],
		) (logr.Logger, error) {
			return l.Value().V(1), nil
		},
	)
}
