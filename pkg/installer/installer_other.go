//go:build !linux && !windows

package installer

import (
	"errors"
	"runtime"
)

// unsupportedInstaller is a no-op installer returned on platforms that are not
// yet supported (e.g. macOS, FreeBSD).
type unsupportedInstaller struct{}

func newPlatform() (Installer, error) {
	return nil, errors.New("installer: service installation is not supported on " + runtime.GOOS +
		"; supported platforms are linux (systemd) and windows (SCM)")
}

func (u *unsupportedInstaller) Install(_ Config) error {
	return errors.New("installer: not supported on " + runtime.GOOS)
}

func (u *unsupportedInstaller) Uninstall(_ Config) error {
	return errors.New("installer: not supported on " + runtime.GOOS)
}

func (u *unsupportedInstaller) Status(_ Config) (string, error) {
	return "", errors.New("installer: not supported on " + runtime.GOOS)
}
