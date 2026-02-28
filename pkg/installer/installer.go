// Package installer provides cross-platform OS service installation (systemd on
// Linux, Service Control Manager on Windows).  Call New() to get a platform-
// specific Installer, then call Install, Uninstall or Status.
package installer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Config holds all parameters needed to register the binary as a service.
type Config struct {
	// Name is the service identifier used by the OS (e.g. "my-portfolio").
	Name string
	// DisplayName is the human-readable label shown in service managers.
	DisplayName string
	// Description is a one-line description of what the service does.
	Description string
	// BinaryPath is the absolute path to the compiled binary.
	// When empty, the running executable path is resolved automatically.
	BinaryPath string
	// WorkDir is the working directory for the service process.
	// When empty, the directory that contains the binary is used.
	WorkDir string
	// User is the OS user the service runs as (Linux only).
	// Defaults to the service Name when empty.
	User string
	// Args are extra command-line arguments appended to ExecStart / binPath.
	Args []string
	// EnvVars are additional KEY=VALUE pairs injected into the service environment.
	EnvVars []string
}

// Installer defines the lifecycle operations for an OS service.
type Installer interface {
	// Install registers, enables and starts the service.
	Install(cfg Config) error
	// Uninstall stops, disables and removes the service registration.
	Uninstall(cfg Config) error
	// Status returns a human-readable description of the service state.
	Status(cfg Config) (string, error)
}

// New returns a platform-specific Installer.
// It returns an error on platforms that are not supported.
func New() (Installer, error) {
	return newPlatform()
}

// Execute is a convenience wrapper that parses the action string ("install",
// "uninstall", "status") and invokes the matching method.
func Execute(action string, cfg Config) error {
	// Resolve the binary path when not provided.
	if cfg.BinaryPath == "" {
		self, err := os.Executable()
		if err != nil {
			return fmt.Errorf("installer: cannot resolve executable path: %w", err)
		}
		self, err = filepath.EvalSymlinks(self)
		if err != nil {
			return fmt.Errorf("installer: cannot evaluate symlinks: %w", err)
		}
		cfg.BinaryPath = self
	}

	// Default working directory to the binary's parent directory.
	if cfg.WorkDir == "" {
		cfg.WorkDir = filepath.Dir(cfg.BinaryPath)
	}

	// Default service user to the service name (Linux).
	if cfg.User == "" {
		cfg.User = cfg.Name
	}

	ins, err := New()
	if err != nil {
		return err
	}

	switch action {
	case "install":
		return ins.Install(cfg)
	case "uninstall":
		return ins.Uninstall(cfg)
	case "status":
		s, err := ins.Status(cfg)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	default:
		return errors.New("installer: unknown action " + action + "; use install, uninstall or status")
	}
}

// run is a small helper that executes a command, printing its combined output,
// and returns a wrapped error on failure.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installer: %s %v: %w", name, args, err)
	}
	return nil
}

// output is like run but captures stdout and returns it as a string.
func output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("installer: %s %v: %w\n%s", name, args, err, out)
	}
	return string(out), nil
}
