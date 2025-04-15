package project

import (
	"fmt"
	"os"
	"path/filepath"
)

func DetermineProjectRoot(currentDir string) (string, error) {
	// Check if the current directory contains flamigo.conf.yml
	for currentDir != "/" {
		configPath := filepath.Join(currentDir, "flamigo.conf.yml")
		if _, err := os.Stat(configPath); err == nil {
			return currentDir, nil
		}
		currentDir = getParentDirectory(currentDir)
	}

	return "", fmt.Errorf("not a Flamigo project")
}

func getParentDirectory(dir string) string {
	return filepath.Dir(dir)
}
