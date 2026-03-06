// Package find provides utility functions to locate valid directories on the filesystem. It is used to find the first existing directory from a list of potential paths, which can be useful for configuration or resource management in applications.
package find

import (
	"fmt"
	"os"
)

// FindValidDirectory iterates through the provided paths and returns the first one that exists on the filesystem. If none of the paths exist, it returns an error.
func FindValidDirectory(paths []string) (string, error) {
	for _, dir := range paths {
		if _, err := os.Stat(dir); err == nil {
			return dir, nil
		}
	}
	return "", fmt.Errorf("no valid report directory found in: %v", paths)
}
