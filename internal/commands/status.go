package commands

import (
	"fmt"

	"github.com/markwharris77/atlas/internal/log"
	"github.com/markwharris77/atlas/internal/runtime"
	"github.com/markwharris77/atlas/internal/tools"
)

type StatusOptions struct {
	Target   tools.Target
	UnitName string
}

func PrintStatus(opts StatusOptions) error {

	if opts.UnitName == "" {
		return fmt.Errorf("name is required")
	}

	status, err := runtime.Status(opts.Target, opts.UnitName)

	if err != nil {
		return fmt.Errorf("failed to get status %w", err)
	}
	log.Info(status)
	return nil
}
