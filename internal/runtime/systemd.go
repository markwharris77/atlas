package runtime

import (
	"fmt"
	"path"
	"strings"

	"github.com/markwharris77/atlas/internal/log"
	"github.com/markwharris77/atlas/internal/tools"
)

type UnitOptions struct {
	Name        string
	Description string
	User        string
	Group       string
	WorkDir     string
	Command     []string
	Env         map[string]string
	Restart     string
	RestartSec  int
}

func RenderUnit(opts UnitOptions) (string, error) {
	if opts.Name == "" {
		return "", fmt.Errorf("systemd unit name is required")
	}
	if len(opts.Command) == 0 {
		return "", fmt.Errorf("systemd command is required")
	}

	var b strings.Builder

	// --- [Unit] ---
	b.WriteString("[Unit]\n")
	if opts.Description != "" {
		b.WriteString("Description=" + opts.Description + "\n")
	} else {
		b.WriteString("Description=" + opts.Name + "\n")
	}
	b.WriteString("After=network.target\n\n")

	// --- [Service] ---
	b.WriteString("[Service]\n")
	b.WriteString("Type=simple\n")

	if opts.User != "" {
		b.WriteString("User=" + opts.User + "\n")
	}
	if opts.Group != "" {
		b.WriteString("Group=" + opts.Group + "\n")
	}

	if opts.WorkDir != "" {
		b.WriteString("WorkingDirectory=" + opts.WorkDir + "\n")
	}

	// Environment variables
	for k, v := range opts.Env {
		b.WriteString(fmt.Sprintf("Environment=%s=%q\n", k, v))
	}

	// ExecStart (argv-style, no shell)
	b.WriteString("ExecStart=" + strings.Join(opts.Command, " ") + "\n")

	// Restart policy
	if opts.Restart != "" {
		b.WriteString("Restart=" + opts.Restart + "\n")
	}
	if opts.RestartSec > 0 {
		b.WriteString(fmt.Sprintf("RestartSec=%d\n", opts.RestartSec))
	}

	b.WriteString("\n")

	// --- [Install] ---
	b.WriteString("[Install]\n")
	b.WriteString("WantedBy=multi-user.target\n")

	return b.String(), nil
}

func InstallUnit(target tools.Target, unitName string, unitContents string, baseDir string) error {
	log.Info("installing systemd unit: %s", unitName)

	if unitName == "" {
		return fmt.Errorf("unit name required")
	}

	UnitFileName := "atlas-" + unitName + ".service"

	remoteTmp := path.Join(baseDir, UnitFileName)
	remoteFinal := path.Join("/etc/systemd/system", UnitFileName)

	cmd := fmt.Sprintf("cat > %q <<'EOF'\n%s\nEOF", remoteTmp, unitContents)

	_, err := tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: cmd,
	})

	if err != nil {
		return fmt.Errorf("write unit file failed: %w", err)
	}

	_, err = tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: fmt.Sprintf("sudo cp %q %q", remoteTmp, remoteFinal),
	})

	if err != nil {
		return fmt.Errorf("install unit failed: %w", err)
	}

	return nil
}

func Reload(target tools.Target) error {
	log.Info("reloading systemd dameon on %s", target.Host)
	_, err := tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: "sudo systemctl daemon-reload",
	})

	if err != nil {
		return fmt.Errorf("reload systemd daemon failed on %s: %w", target.Host, err)
	}

	return nil
}

func Restart(target tools.Target, unitName string) error {

	log.Info("restarting systemd unit %s on %s", unitName, target.Host)

	if unitName == "" {
		return fmt.Errorf("name required")
	}

	_, err := tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: fmt.Sprintf("sudo systemctl restart %s", "atlas-"+unitName),
	})

	if err != nil {
		return fmt.Errorf("restart unit %s failed: %w", unitName, err)
	}

	return nil
}

func Status(target tools.Target, unitName string) (string, error) {

	log.Info("getting systemd status for unit %s on %s", unitName, target.Host)

	if unitName == "" {
		return "unknown", fmt.Errorf("name required")
	}

	out, err := tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: fmt.Sprintf("sudo systemctl status %s || true", "atlas-"+unitName),
	})

	if err != nil {
		return "unkown", fmt.Errorf("failed to get status for unit %s: %w", unitName, err)
	}

	return out, nil
}

func Enable(target tools.Target, unitName string) error {
	log.Info("enabling systemd unit %s on %s", unitName, target.Host)

	if unitName == "" {
		return fmt.Errorf("name required")
	}

	_, err := tools.RunSSH(tools.SSHOptions{
		Target:  target,
		Command: fmt.Sprintf("sudo systemctl enable %s", "atlas-"+unitName),
	})

	if err != nil {
		return fmt.Errorf("failed to enable unit %s: %w", unitName, err)
	}

	return nil
}