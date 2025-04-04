package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var defaultDirectories = []string{
	"{root}",
	"{root}/cmd",
	"{root}/internal/",
	"{root}/internal/domains",
	"{root}/internal/api",
	"{root}/internal",
	"{root}/interfaces",
}

func InitializeDirectories(rootPath string) error {
	for _, dir := range defaultDirectories {
		dirPath := strings.ReplaceAll(dir, "{root}", rootPath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}
	return nil
}

func InitializeGoMod(rootPath string, moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = rootPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("initialize go module: %w", err)
	}
	return nil
}

func TidyGoMod(rootPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = rootPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tidy go module: %w", err)
	}
	return nil
}
