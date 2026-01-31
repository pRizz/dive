package adapter

import (
	"context"
	"fmt"
	"github.com/pRizz/dive/cmd/dive/cli/internal/command/export"
	"github.com/pRizz/dive/dive/image"
	"github.com/pRizz/dive/internal/bus"
	"github.com/pRizz/dive/internal/bus/event/payload"
	"github.com/pRizz/dive/internal/log"
	"github.com/spf13/afero"
	"os"
)

type Exporter interface {
	ExportTo(ctx context.Context, img *image.Analysis, path string) error
}

type jsonExporter struct {
	filesystem afero.Fs
}

func NewExporter(fs afero.Fs) Exporter {
	return &jsonExporter{
		filesystem: fs,
	}
}

func (e *jsonExporter) ExportTo(ctx context.Context, analysis *image.Analysis, path string) error {
	log.WithFields("path", path).Infof("exporting analysis")

	mon := bus.StartTask(payload.GenericTask{
		Title: payload.Title{
			Default:      "Exporting details",
			WhileRunning: "Exporting details",
			OnSuccess:    "Exported details",
		},
		HideOnSuccess:      false,
		HideStageOnSuccess: false,
		ID:                 analysis.Image,
		Context:            fmt.Sprintf("[file: %s]", path),
	})

	bytes, err := export.NewExport(analysis).Marshal()
	if err != nil {
		mon.SetError(err)
		return fmt.Errorf("cannot marshal export payload: %w", err)
	} else {
		mon.SetCompleted()
	}

	file, err := e.filesystem.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open export file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.WithFields("error", closeErr, "path", path).Warn("failed to close export file")
		}
	}()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("cannot write to export file: %w", err)
	}
	return nil
}
