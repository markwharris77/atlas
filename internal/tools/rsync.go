package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markwharris77/atlas/internal/log"
)

type RsyncOptions struct {
	LocalPath  string // e.g. "./cmd"
	RemotePath string // e.g. "/home/mark/test_service"
	Target     Target
	Delete     bool // delete remote files not present locally
	DryRun     bool
}

func RunRsync(opts RsyncOptions) error {

	log.Info("syncing %s into %s on %s", opts.LocalPath, opts.RemotePath, opts.Target.Host)

	if opts.LocalPath == "" || opts.RemotePath == "" || opts.Target.Addr() == "" {
		return fmt.Errorf("rsync: LocalPath, RemotePath, and Host are required")
	}

	local := opts.LocalPath
	if fi, err := filepath.Abs(local); err == nil {
		local = fi
	}
	if !strings.HasSuffix(local, string(filepath.Separator)) {
		local += string(filepath.Separator)
	}

	remoteTarget := fmt.Sprintf("%s:%s", opts.Target.Addr(), opts.RemotePath)

	args := []string{
		"-az", // archive + compress
	}

	if opts.Delete {
		args = append(args, "--delete")
	}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}

	args = append(args, local, remoteTarget)

	cmd := exec.Command("rsync", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	log.Verbose("running: rsync %s", strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsync failed: %w\noutput:\n%s", err, out.String())
	}

	if s := strings.TrimSpace(out.String()); s != "" {
		log.Verbose("rsync output: %s", s)
	}
	return nil
}
