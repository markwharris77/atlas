package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/markwharris77/atlas/internal/log"
)

type SSHOptions struct {
	Target  Target
	Command string
	Args    []string
}

func RunSSH(opts SSHOptions) (string, error) {
	if opts.Target.Host == "" {
		return "", fmt.Errorf("ssh: target host is required")
	}

	sshArgs := []string{
		"-o", "BatchMode=yes",
	}

	if opts.Target.Port != 0 {
		sshArgs = append(sshArgs, "-p", strconv.Itoa(opts.Target.Port))
	}

	sshArgs = append(sshArgs, opts.Target.Addr(), "--")

	remoteCmd := opts.Command
	if remoteCmd == "" && len(opts.Args) > 0 {
		remoteCmd = strings.Join(opts.Args, " ")
	}
	if strings.TrimSpace(remoteCmd) == "" {
		return "", fmt.Errorf("ssh: command is required")
	}

	sshArgs = append(sshArgs, remoteCmd)

	log.Verbose("running: ssh %s", strings.Join(sshArgs, " "))

	cmd := exec.Command("ssh", sshArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := out.String()
	if err != nil {
		return output, fmt.Errorf("ssh failed: %w\noutput:\n%s", err, output)
	}
	return output, nil
}
