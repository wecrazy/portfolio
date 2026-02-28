//go:build windows

package installer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// windowsInstaller manages the service via NSSM (Non-Sucking Service Manager).
// https://nssm.cc
// Must be run from an elevated (Administrator) prompt.
type windowsInstaller struct {
	// nssmPath is the resolved absolute path to nssm.exe.
	nssmPath string
}

// wellKnownNSSMPaths lists common installation directories to probe when nssm
// is not on the system PATH.
var wellKnownNSSMPaths = []string{
	`C:\nssm\nssm.exe`,
	`C:\nssm\win64\nssm.exe`,
	`C:\nssm\win32\nssm.exe`,
	`C:\Program Files\nssm\nssm.exe`,
	`C:\Program Files\nssm\win64\nssm.exe`,
	`C:\Program Files (x86)\nssm\nssm.exe`,
	`C:\tools\nssm\nssm.exe`,
	`C:\tools\nssm.exe`,
	`C:\ProgramData\chocolatey\bin\nssm.exe`, // Chocolatey install
	`C:\Windows\nssm.exe`,
}

// newPlatform resolves nssm.exe — from PATH, well-known locations, or user
// prompt — and returns a ready-to-use installer.
func newPlatform() (Installer, error) {
	nssmPath, err := resolveNSSM()
	if err != nil {
		return nil, err
	}
	return &windowsInstaller{nssmPath: nssmPath}, nil
}

// resolveNSSM tries to find nssm.exe and, if it cannot, asks the user for the
// path.  It returns an error only when no usable path is provided.
func resolveNSSM() (string, error) {
	// 1. Check system PATH first.
	for _, name := range []string{"nssm.exe", "nssm"} {
		if p, err := exec.LookPath(name); err == nil {
			fmt.Printf("[installer] Found NSSM on PATH: %s\n", p)
			return p, nil
		}
	}

	// 2. Probe well-known installation paths.
	for _, p := range wellKnownNSSMPaths {
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("[installer] Found NSSM at: %s\n", p)
			return p, nil
		}
	}

	// 3. NSSM not found — prompt the user.
	fmt.Println()
	fmt.Println("[installer] NSSM (Non-Sucking Service Manager) was not found on this system.")
	fmt.Println("            NSSM is required to install the service on Windows.")
	fmt.Println("            Download it from: https://nssm.cc/download")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter the full path to nssm.exe (or 'q' to cancel): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("installer: failed to read input: %w", err)
		}

		p := strings.TrimSpace(line)
		if p == "q" || p == "Q" {
			return "", fmt.Errorf("installer: installation cancelled by user")
		}
		if p == "" {
			fmt.Println("  Path cannot be empty.")
			continue
		}

		// Validate the given path.
		if _, err := os.Stat(p); err != nil {
			fmt.Printf("  File not found: %s\n", p)
			continue
		}
		fmt.Printf("[installer] Using NSSM at: %s\n", p)
		return p, nil
	}
}

// nssmSetting is a key/value pair passed to `nssm set`.
type nssmSetting struct {
	key   string
	value string
}

// nssm is a helper that runs the resolved nssm.exe with the given arguments.
func (w *windowsInstaller) nssm(args ...string) error {
	return run(w.nssmPath, args...)
}

// nssmOut is like nssm but captures and returns combined output.
func (w *windowsInstaller) nssmOut(args ...string) (string, error) {
	return output(w.nssmPath, args...)
}

// Install registers the binary as a Windows service using NSSM and starts it.
func (w *windowsInstaller) Install(cfg Config) error {
	fmt.Printf("[installer] Installing service %q via NSSM...\n", cfg.Name)

	// nssm install <name> <program> [args...]
	installArgs := []string{"install", cfg.Name, cfg.BinaryPath}
	installArgs = append(installArgs, cfg.Args...)
	if err := w.nssm(installArgs...); err != nil {
		return fmt.Errorf("installer: nssm install failed: %w", err)
	}

	// Configure service properties.
	settings := []nssmSetting{
		{key: "AppDirectory", value: cfg.WorkDir},
		{key: "DisplayName", value: cfg.DisplayName},
		{key: "Description", value: cfg.Description},
		{key: "Start", value: "SERVICE_AUTO_START"},
		// Route stdout/stderr to a log file in the working directory.
		{key: "AppStdout", value: cfg.WorkDir + `\logs\` + cfg.Name + `-stdout.log`},
		{key: "AppStderr", value: cfg.WorkDir + `\logs\` + cfg.Name + `-stderr.log`},
		{key: "AppRotateFiles", value: "1"},
		{key: "AppRotateSeconds", value: "86400"},
		{key: "AppRotateBytes", value: "10485760"}, // 10 MB
	}

	for _, s := range settings {
		if s.value == "" {
			continue
		}
		fmt.Printf("[installer] nssm set %s %s %s\n", cfg.Name, s.key, s.value)
		if err := w.nssm("set", cfg.Name, s.key, s.value); err != nil {
			// Non-fatal: log and continue.
			fmt.Printf("[installer] Warning: could not set %s: %v\n", s.key, err)
		}
	}

	// Inject extra environment variables.
	for _, env := range cfg.EnvVars {
		fmt.Printf("[installer] nssm set %s AppEnvironmentExtra %s\n", cfg.Name, env)
		if err := w.nssm("set", cfg.Name, "AppEnvironmentExtra", env); err != nil {
			fmt.Printf("[installer] Warning: could not set env %s: %v\n", env, err)
		}
	}

	// Create log directory so NSSM can write to it immediately.
	logDir := cfg.WorkDir + `\logs`
	_ = os.MkdirAll(logDir, 0o755)

	// Start the service.
	fmt.Printf("[installer] Starting service %q...\n", cfg.Name)
	if err := w.nssm("start", cfg.Name); err != nil {
		return fmt.Errorf("installer: failed to start service: %w", err)
	}

	fmt.Printf("\n[installer] Service %q installed and started.\n", cfg.Name)
	fmt.Printf("  nssm status %s\n", cfg.Name)
	fmt.Printf("  nssm edit   %s          (GUI editor)\n", cfg.Name)
	fmt.Printf("  nssm restart %s\n", cfg.Name)
	fmt.Printf("  nssm stop   %s\n", cfg.Name)
	fmt.Printf("  Logs: %s\n", logDir)
	return nil
}

// Uninstall stops and removes the Windows service via NSSM.
func (w *windowsInstaller) Uninstall(cfg Config) error {
	fmt.Printf("[installer] Stopping service %q...\n", cfg.Name)
	_ = w.nssm("stop", cfg.Name) // ignore error — may already be stopped

	fmt.Printf("[installer] Removing service %q...\n", cfg.Name)
	// nssm remove <name> confirm  (the "confirm" flag suppresses the interactive prompt)
	if err := w.nssm("remove", cfg.Name, "confirm"); err != nil {
		return fmt.Errorf("installer: nssm remove failed: %w", err)
	}

	fmt.Printf("[installer] Service %q uninstalled.\n", cfg.Name)
	return nil
}

// Status returns the output of `nssm status <name>`.
func (w *windowsInstaller) Status(cfg Config) (string, error) {
	return w.nssmOut("status", cfg.Name)
}
