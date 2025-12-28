package commands

import (
	"fmt"

	"github.com/markwharris77/atlas/internal/log"
	"github.com/markwharris77/atlas/internal/runtime"
	"github.com/markwharris77/atlas/internal/tools"
)

type RestartOptions struct {
	Target   tools.Target
	UnitName string
}

func RunRestart(opts RestartOptions) error {
	log.Info("restarting unit %s on %s", opts.UnitName, opts.Target.Host)
	err := runtime.Restart(opts.Target, opts.UnitName)

	if err != nil {
		return fmt.Errorf("failed to restart unit %s: %w", opts.UnitName, err)
	}

	status, err := runtime.Status(opts.Target, opts.UnitName)

	if err != nil {
		return fmt.Errorf("failed to get status after restart: %w", err)
	}
	log.Info("Status:%s\n", status)
	return nil
}
