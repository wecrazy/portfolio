//go:build linux

package installer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// linuxInstaller manages the service via systemd.
type linuxInstaller struct{}

func newPlatform() (Installer, error) {
	// Verify systemd is available.
	if _, err := exec.LookPath("systemctl"); err != nil {
		return nil, fmt.Errorf("installer: systemctl not found; only systemd-based Linux distributions are supported")
	}
	return &linuxInstaller{}, nil
}

// serviceFile returns the absolute path of the systemd unit file.
func (*linuxInstaller) serviceFile(name string) string {
	return fmt.Sprintf("/etc/systemd/system/%s.service", name)
}

// unitTemplate is the systemd service unit content.
var unitTemplate = template.Must(template.New("unit").Parse(`[Unit]
Description={{ .Description }}
After=network.target
Wants=network.target

[Service]
Type=simple
User={{ .User }}
Group={{ .User }}
WorkingDirectory={{ .WorkDir }}
ExecStart={{ .ExecStart }}
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier={{ .Name }}
{{ range .EnvVars }}
Environment={{ . }}
{{ end }}
# Security hardening
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths={{ .WorkDir }} {{ .WorkDir }}/data {{ .WorkDir }}/uploads {{ .WorkDir }}/logs

[Install]
WantedBy=multi-user.target
`))

type unitData struct {
	Name        string
	Description string
	User        string
	WorkDir     string
	ExecStart   string
	EnvVars     []string
}

// Install creates the systemd unit file, reloads the daemon, enables and
// starts the service.  Must be run as root.
func (l *linuxInstaller) Install(cfg Config) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("installer: must be run as root (sudo ./bin/%s --install)", cfg.Name)
	}

	// Check if the service is already running.
	if isActive(cfg.Name) {
		fmt.Printf("[installer] Service %q is already running. Restart it with:\n  sudo systemctl restart %s\n", cfg.Name, cfg.Name)
		return nil
	}

	// Ensure the service user exists.
	if err := ensureUser(cfg.User); err != nil {
		return err
	}

	// Build ExecStart string.
	execStart := cfg.BinaryPath
	if len(cfg.Args) > 0 {
		execStart += " " + strings.Join(cfg.Args, " ")
	}

	// Render the unit file.
	var buf bytes.Buffer
	data := unitData{
		Name:        cfg.Name,
		Description: cfg.Description,
		User:        cfg.User,
		WorkDir:     cfg.WorkDir,
		ExecStart:   execStart,
		EnvVars:     cfg.EnvVars,
	}
	if err := unitTemplate.Execute(&buf, data); err != nil {
		return fmt.Errorf("installer: render unit file: %w", err)
	}

	// Write the unit file.
	unitPath := l.serviceFile(cfg.Name)
	fmt.Printf("[installer] Writing unit file: %s\n", unitPath)
	if err := os.WriteFile(unitPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("installer: write unit file %s: %w", unitPath, err)
	}

	// Ownership of WorkDir.
	_ = run("chown", "-R", cfg.User+":"+cfg.User, cfg.WorkDir)

	// Reload, enable and start.
	steps := [][]string{
		{"systemctl", "daemon-reload"},
		{"systemctl", "enable", cfg.Name},
		{"systemctl", "start", cfg.Name},
	}
	for _, s := range steps {
		fmt.Printf("[installer] Running: %s\n", strings.Join(s, " "))
		if err := run(s[0], s[1:]...); err != nil {
			return err
		}
	}

	fmt.Printf("\n[installer] Service %q installed and started.\n", cfg.Name)
	fmt.Printf("  sudo systemctl status %s\n", cfg.Name)
	fmt.Printf("  sudo journalctl -u %s -f\n", cfg.Name)
	return nil
}

// Uninstall stops, disables and removes the service unit file.
func (l *linuxInstaller) Uninstall(cfg Config) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("installer: must be run as root (sudo ./bin/%s --uninstall)", cfg.Name)
	}

	steps := [][]string{
		{"systemctl", "stop", cfg.Name},
		{"systemctl", "disable", cfg.Name},
	}
	for _, s := range steps {
		fmt.Printf("[installer] Running: %s\n", strings.Join(s, " "))
		// Ignore errors — the service may already be stopped/disabled.
		_ = run(s[0], s[1:]...)
	}

	unitPath := l.serviceFile(cfg.Name)
	fmt.Printf("[installer] Removing unit file: %s\n", unitPath)
	if err := os.Remove(unitPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("installer: remove unit file: %w", err)
	}

	_ = run("systemctl", "daemon-reload")

	fmt.Printf("[installer] Service %q uninstalled.\n", cfg.Name)
	return nil
}

// Status returns the output of `systemctl status <name>`.
func (*linuxInstaller) Status(cfg Config) (string, error) {
	out, err := exec.Command("systemctl", "status", cfg.Name, "--no-pager").CombinedOutput()
	if err != nil {
		// systemctl status exits with code 3 when inactive — still useful output.
		return string(out), nil
	}
	return string(out), nil
}

// isActive reports whether the service is currently active.
func isActive(name string) bool {
	return exec.Command("systemctl", "is-active", "--quiet", name).Run() == nil
}

// ensureUser creates a system user if it does not already exist.
func ensureUser(user string) error {
	if exec.Command("id", user).Run() == nil {
		fmt.Printf("[installer] Service user %q already exists.\n", user)
		return nil
	}
	fmt.Printf("[installer] Creating system user %q...\n", user)
	return run("useradd", "--system", "--shell", "/bin/false",
		"--no-create-home", user)
}
