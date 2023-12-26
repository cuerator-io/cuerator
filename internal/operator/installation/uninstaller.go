package installation

import (
	"context"

	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	"github.com/go-logr/logr"
)

type uninstaller struct {
	Installation *model.Installation
	Logger       logr.Logger
}

func (u *uninstaller) Run(context.Context) error {
	return nil
}
