package commands

import (
	"fmt"
	"path"
	"time"

	"github.com/markwharris77/atlas/internal/log"
	"github.com/markwharris77/atlas/internal/runtime"
	"github.com/markwharris77/atlas/internal/tools"
)

type DeployOptions struct {
	Name       string
	LocalPath  string
	Target     tools.Target
	RunUser    string
	RunGroup   string
	RunCommand []string
	RunEnv     map[string]string
}

func RunDeploy(opts DeployOptions) error {

	log.Info("deploying %s to %s", opts.Name, opts.Target.Host)

	basePath := fmt.Sprintf("/opt/atlas/%s/", opts.Name)
	releases := path.Join(basePath, "releases")
	releaseID := time.Now().Format("20060102-150405")
	releasePath := path.Join(releases, releaseID)
	currentLink := path.Join(basePath, "current")

	log.Info("release: %s", releaseID)

	_, err := tools.RunSSH(tools.SSHOptions{
		Target:  opts.Target,
		Command: fmt.Sprintf("mkdir -p %s", releases),
	})

	if err != nil {
		return fmt.Errorf("deploy: prepare dirs: %w", err)
	}

	err = tools.RunRsync(tools.RsyncOptions{
		LocalPath:  opts.LocalPath,
		RemotePath: releasePath,
		Target:     opts.Target,
		Delete:     true,
		DryRun:     false,
	})

	if err != nil {
		return fmt.Errorf("deploy: rsync: %w\n", err)
	}

	log.Info("pointing %s to %s", currentLink, releasePath)

	switchCmd := fmt.Sprintf(
		"ln -sfn %q %q && mv -Tf %q %q",
		releasePath, path.Join(basePath, ".current.tmp"),
		path.Join(basePath, ".current.tmp"), currentLink,
	)

	_, err = tools.RunSSH(tools.SSHOptions{
		Target:  opts.Target,
		Command: switchCmd,
	})
	if err != nil {
		return fmt.Errorf("deploy: switch current: %w", err)
	}

	unit, err := runtime.RenderUnit(runtime.UnitOptions{
		Name:        opts.Name,
		Description: "applicaiton deployed via atlas",
		User:        opts.RunUser,
		Group:       opts.RunGroup,
		Command:     opts.RunCommand,
		Env:         opts.RunEnv,
		Restart:     "always",
		RestartSec:  30,
	})

	if err != nil {
		return fmt.Errorf("deploy: create unit: %w", err)
	}

	err = runtime.InstallUnit(
		opts.Target,
		opts.Name,
		unit,
		releasePath,
	)

	if err != nil {
		return fmt.Errorf("deploy: install unit: %w", err)
	}

	err = runtime.Reload(opts.Target)

	if err != nil {
		return fmt.Errorf("deploy: failed to restart daemon: %w", err)
	}

	err = runtime.Enable(opts.Target, opts.Name)

	if err != nil {
		return fmt.Errorf("deploy: failed to enable unit: %w", err)
	}

	err = runtime.Restart(opts.Target, opts.Name)

	if err != nil {
		return fmt.Errorf("deploy: failed to restart unit: %w", err)
	}

	status, err := runtime.Status(opts.Target, opts.Name)

	if err != nil {
		log.Info("failed to get status: %w", err)
	}

	log.Info("status:\n%s", status)

	return nil
}
