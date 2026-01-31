package view

import (
	"github.com/pRizz/dive/cmd/dive/cli/internal/ui/v1/viewmodel"
)

type LayerChangeListener func(viewmodel.LayerSelection) error
